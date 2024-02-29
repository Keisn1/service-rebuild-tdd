package main

import (
	"github.com/Keisn1/note-taking-app/server"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type config struct {
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
	handler := http.HandlerFunc(server.NotesService)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Address, handler))
}
