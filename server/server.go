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
	AddNote(userID int, note string) error
}

type NotesServer struct {
	NotesStore NotesStore
}

func (ns *NotesServer) ProcessAddNote(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/notes/"))
	if err != nil {
		log.Println(fmt.Errorf("NotesServer.ServeHTTP: %w", err))
		http.Error(w, "There was an Error retrieving Notes", http.StatusInternalServerError)
	}

	note := "sample note"
	err = ns.NotesStore.AddNote(userID, note)
	w.WriteHeader(http.StatusAccepted)

}

func (ns *NotesServer) GetNotes(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/notes" {
		notes := ns.NotesStore.GetAllNotes()
		json.NewEncoder(w).Encode(notes)
		return
	} else {
		userID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/notes/"))
		if err != nil {
			log.Println(fmt.Errorf("NotesServer.ServeHTTP: %w", err))
			http.Error(w, "There was an Error retrieving Notes", http.StatusInternalServerError)
		}
		notes := ns.NotesStore.GetNotesByID(userID)
		if len(notes) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(notes)
		return
	}
}

func (ns *NotesServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ns.ProcessAddNote(w, r)
	case http.MethodGet:
		ns.GetNotes(w, r)
	}
}
