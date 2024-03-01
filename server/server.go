package server

import (
	"encoding/json"
	"net/http"
)

type NotesStore interface {
	GetAllNotes() []string
}

type NotesServer struct {
	NotesStore NotesStore
}

func (ns *NotesServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	notes := ns.NotesStore.GetAllNotes()
	json.NewEncoder(w).Encode(notes)
}
