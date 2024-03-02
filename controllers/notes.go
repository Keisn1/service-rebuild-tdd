package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Note struct {
	UserID int
	Note   string
}

func NewNote(userID int, note string) Note {
	return Note{userID, note}
}

type Notes []Note

type NotesStore interface {
	GetAllNotes() Notes
	GetNotesByID(int) Notes
	AddNote(Note) error
}

type NotesCtrlr struct {
	NotesStore NotesStore
}

func NewNotesController(store NotesStore) NotesCtrlr {
	return NotesCtrlr{store}
}

func (ns *NotesCtrlr) ProcessAddNote(w http.ResponseWriter, r *http.Request) {
	var body map[string]Note
	_ = json.NewDecoder(r.Body).Decode(&body)

	note, _ := body["note"]
	_ = ns.NotesStore.AddNote(note)
	w.WriteHeader(http.StatusAccepted)
}

func (ns *NotesCtrlr) GetNotesByID(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	notes := ns.NotesStore.GetNotesByID(userID)
	if len(notes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode([]string{})
		return
	}
	json.NewEncoder(w).Encode(notes)
	return
}

func (ns *NotesCtrlr) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes := ns.NotesStore.GetAllNotes()
	json.NewEncoder(w).Encode(notes)
	return
}
