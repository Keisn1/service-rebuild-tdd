package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type NotesStore interface {
	GetAllNotes() []string
	GetNotesByID(int) []string
}

type NotesServer struct {
	NotesStore NotesStore
}

func (ns *NotesServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/notes" {
		notes := ns.NotesStore.GetAllNotes()
		json.NewEncoder(w).Encode(notes)
		return
	} else {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/notes/"))
		if err != nil {
			log.Println(fmt.Errorf("NotesServer.ServeHTTP: %w", err))
			http.Error(w, "There was an Error retrieving Notes", http.StatusInternalServerError)
		}
		notes := ns.NotesStore.GetNotesByID(id)
		json.NewEncoder(w).Encode(notes)
		return
	}
}
