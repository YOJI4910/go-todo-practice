package main

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormDB *gorm.DB

type TODO struct {
	gorm.Model
	Author      string    `json:"author"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
}

func list(response http.ResponseWriter, request *http.Request) {
	todos := make([]TODO, 0)
	if err := gormDB.Find(&todos).Error; err != nil {
		panic(err)
	}

	json.NewEncoder(response).Encode(todos)
}

func add(response http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(response)
	if request.Body == http.NoBody {
		_ = encoder.Encode("request body is nil")
		return
	}

	var todo TODO
	if err := json.NewDecoder(request.Body).Decode(&todo); err != nil {
		_ = encoder.Encode(err.Error())
		return
	}

	if err := gormDB.Create(&todo).Error; err != nil {
		_ = encoder.Encode(err.Error())
		return
	}

	encoder.Encode(todo)
}

func main() {
	_gormDB, err := gorm.Open(mysql.Open("root:mysql@tcp(goTodoDB)/goTodo?parseTime=true"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db, err := _gormDB.DB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if !_gormDB.Migrator().HasTable(&TODO{}) {
		if err := _gormDB.Migrator().CreateTable(&TODO{}); err != nil {
			panic(err)
		}
	}

	gormDB = _gormDB

	http.HandleFunc("/todo/list", list)
	http.HandleFunc("/todo/add", add)
	_ = http.ListenAndServe(":8080", nil)
}
