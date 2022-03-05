package main

import (
	"encoding/json"
	"net/http"
	"time"
)

var todos []TODOObject = []TODOObject{}

type MetaData struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type TODO struct {
	Author      string    `json:"author"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
}

type TODOObject struct {
	MetaData
	TODO
}

func list(response http.ResponseWriter, request *http.Request) {
	json.NewEncoder(response).Encode(todos)
}

func add(response http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(response)
	if request.Body == http.NoBody {
		encoder.Encode("request body is nil")
		return
	}
	var todo TODO
	if err := json.NewDecoder(request.Body).Decode(&todo); err != nil {
		encoder.Encode(err.Error())
		return
	}
	metadata := MetaData{
		ID:        len(todos) + 1,
		CreatedAt: time.Now(),
	}
	object := TODOObject{
		TODO:     todo,
		MetaData: metadata,
	}
	todos = append(todos, object)
	encoder.Encode(object)
}

func main() {
	http.HandleFunc("/todo/list", list)
	http.HandleFunc("/todo/add", add)
	http.ListenAndServe(":8080", nil)
}
