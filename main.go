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

type InMemoryPlayerStore struct {
	notes map[int][]string
}

func (i *InMemoryPlayerStore) GetNotesByID(id int) []string {
	return i.notes[id]
}

func (i *InMemoryPlayerStore) GetAllNotes() []string {
	var allNotes []string
	for _, notes := range i.notes {
		allNotes = append(allNotes, notes...)
	}
	return allNotes
}

func (i *InMemoryPlayerStore) AddNote(userID int, note string) error {
	i.notes[userID] = append(i.notes[userID], note)
	return nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	notesC := &controllers.Notes{NotesStore: &InMemoryPlayerStore{}}

	r := chi.NewRouter()

	r.Route("/notes", func(r chi.Router) {
		r.Get("/", notesC.GetAllNotes)
		r.Get("/{id}", notesC.GetNotesByID)
		r.Post("/{id}", notesC.ProcessAddNote)
	})

	log.Fatal(http.ListenAndServe(":"+cfg.Server.Address, r))
}
