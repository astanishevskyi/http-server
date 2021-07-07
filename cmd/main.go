package main

import (
	"flag"
	"github.com/astanishevskyi/http-server/internal/apiserver"
	"github.com/astanishevskyi/http-server/internal/apiserver/configs"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/.env", "path to config file")
}

func main() {
	flag.Parse()

	err := godotenv.Load(configPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	grpcAddr := os.Getenv("GRPCADDR")
	config := configs.Config{
		BindAddr: port,
		GRPCAddr: grpcAddr,
	}
	server := apiserver.New(&config)

	server.ConfigRouter()
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
