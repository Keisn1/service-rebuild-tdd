package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Notes []string

type NotesStore interface {
	GetAllNotes() map[int]Notes
	GetNotesByID(int) Notes
	AddNote(userID int, note string) error
}

type NotesCtrlr struct {
	NotesStore NotesStore
}

func NewNotesController(store NotesStore) NotesCtrlr {
	return NotesCtrlr{store}
}

func (ns *NotesCtrlr) ProcessAddNote(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var body map[string]string
	_ = json.NewDecoder(r.Body).Decode(&body)

	note := body["note"]
	_ = ns.NotesStore.AddNote(userID, note)
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
