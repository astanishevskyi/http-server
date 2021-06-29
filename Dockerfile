FROM golang:1.16-alpine

ENV config_path $config_path
RUN mkdir http-server
WORKDIR http-server
COPY . .

RUN go mod tidy
RUN echo $config_path
RUN go build ./cmd/main.go
EXPOSE 8080
