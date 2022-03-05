FROM golang:1.17.8-alpine as dev

ENV ROOT=/go/src/app
# コンパイル時にC言語ライブラリを使うかどうか： ０＝ off
ENV CGO_ENABLED 0
WORKDIR ${ROOT}

RUN apk update && apk add git

# NOTE: Install Linter
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
  sh -s -- -b $(go env GOPATH)/bin v1.35.2

# NOTE: Install hot reload library
RUN go install github.com/cosmtrek/air@v1.27.3

COPY . .

# 依存ライブラリのダウンロード
RUN go mod tidy

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
