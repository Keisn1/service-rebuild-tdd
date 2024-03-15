package main

import (
	"github.com/Keisn1/note-taking-app/controllers"
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

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	notesStore := controllers.NewInMemoryNotesStore()
	notesC := &controllers.NotesCtrlr{NotesStore: notesStore}

	r := chi.NewRouter()

	r.Route("/notes", func(r chi.Router) {
		r.Get("/", notesC.GetAllNotes)
		r.Get("/{id}", notesC.GetNotesByUserID)
		r.Post("/{id}", notesC.ProcessAddNote)
	})

	log.Fatal(http.ListenAndServe(":"+cfg.Server.Address, r))
}
