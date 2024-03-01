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

type InMemoryNotesStore struct {
	notes map[int][]string
}

func (i *InMemoryNotesStore) GetNotesByID(id int) []string {
	return i.notes[id]
}

func (i *InMemoryNotesStore) GetAllNotes() []string {
	var allNotes []string
	for _, notes := range i.notes {
		allNotes = append(allNotes, notes...)
	}
	return allNotes
}

func (i *InMemoryNotesStore) AddNote(userID int, note string) error {
	i.notes[userID] = append(i.notes[userID], note)
	return nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	notesC := &controllers.Notes{NotesStore: &InMemoryNotesStore{
		notes: map[int][]string{
			1: {"Note 1 user 1", "Note 2 user 1"},
			2: {"Note 1 user 2", "Note 2 user 2"},
		},
	}}

	r := chi.NewRouter()

	r.Route("/notes", func(r chi.Router) {
		r.Get("/", notesC.GetAllNotes)
		r.Get("/{id}", notesC.GetNotesByID)
		r.Post("/{id}", notesC.ProcessAddNote)
	})

	log.Fatal(http.ListenAndServe(":"+cfg.Server.Address, r))
}
