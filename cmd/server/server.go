package main

import (
	"github.com/Keisn1/note-taking-app/controllers"
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

	notesStore := controllers.NewInMemoryNotesStore()
	notesC := &controllers.NotesCtrlr{NotesStore: notesStore}

	r := chi.NewRouter()

	r.Route("/users/", func(r chi.Router) {
		r.Get("{userID}/notes", notesC.GetNotesByUserID)
		r.Get("/{userID}/notes/{noteID}", notesC.GetNoteByUserIDAndNoteID)
		r.Post("/{userID}/notes", notesC.Add)
		r.Put("/{userID}/notes/{noteID}", notesC.Edit)
		r.Delete("/{userID}/notes/{noteID}", notesC.Delete)
	})

	log.Fatal(http.ListenAndServe(":"+cfg.Server.Address, r))
}
