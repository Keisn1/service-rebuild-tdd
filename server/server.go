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
	if r.URL.Path == "/notes/1" {
		json.NewEncoder(w).Encode(
			[]string{"Note 1 user 1", "Note 2 user 1"},
		)
	}
	notes := ns.NotesStore.GetAllNotes()
	json.NewEncoder(w).Encode(notes)
}
