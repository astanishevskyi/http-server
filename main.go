package main

import (
	"github.com/astanishevskyi/http-server/apiserver"
	"github.com/astanishevskyi/http-server/apiserver/configs"
	"log"
)

func main() {
	config := configs.Config{
		BindAddr: ":8080",
		Storage:  "in-memory",
	}
	server := apiserver.New(&config)

	if err := server.ConfigStorage(); err != nil {
		log.Fatal(err)
	}
	server.ConfigRouter()
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
