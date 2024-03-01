package main

import (
	"github.com/Keisn1/note-taking-app/server"
	"github.com/go-chi/chi"
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

type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetAllNotes() []string {
	return []string{"Note1", "Note2"}
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	server := &server.NotesServer{NotesStore: &InMemoryPlayerStore{}}

	r := chi.NewRouter()
	r.Get("/notes", server.ServeHTTP)

	log.Fatal(http.ListenAndServe(":"+cfg.Server.Address, r))
}
