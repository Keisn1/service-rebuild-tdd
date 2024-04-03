package main

import (
	"github.com/Keisn1/note-taking-app/app/handlers/notesgrp"
	"github.com/go-chi/chi/v5"
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

	notesStore := notesgrp.NewInMemoryNotesStore()
	hdl := notesgrp.NewHandlers(notesStore)
	r := chi.NewRouter()

	r.Route("/users/", func(r chi.Router) {
		r.Get("/notes", hdl.GetAllNotes)
		r.Get("/{userID}/notes", hdl.GetNotesByUserID)
		r.Get("/{userID}/notes/{noteID}", hdl.GetNoteByUserIDAndNoteID)
		r.Post("/{userID}/notes", hdl.Add)
		r.Put("/{userID}/notes/{noteID}", hdl.Edit)
		r.Delete("/{userID}/notes/{noteID}", hdl.Delete)
	})

	log.Fatal(http.ListenAndServe(":"+cfg.Server.Address, r))
}
