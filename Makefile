all: build run
.PHONY: build
config=
build:
	go build cmd/main.go
run:
	go run cmd/main.go
build-container:
	docker build -t http-server .
up:
	docker run -p 8080:8080 http-server ./main $(config)
