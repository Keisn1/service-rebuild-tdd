package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
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

type Notes struct {
	NotesStore NotesStore
}

func (ns *Notes) ProcessAddNote(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/notes/"))
	if err != nil {
		log.Println(fmt.Errorf("NotesServer.ServeHTTP: %w", err))
		http.Error(w, "There was an Error retrieving Notes", http.StatusInternalServerError)
	}

	note := "sample note"
	err = ns.NotesStore.AddNote(userID, note)
	w.WriteHeader(http.StatusAccepted)
}

func (ns *Notes) GetNotesByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
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

func (ns *Notes) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes := ns.NotesStore.GetAllNotes()
	json.NewEncoder(w).Encode(notes)
	return
}
