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

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type NotesCtrlr struct {
	NotesStore NotesStore
	Logger     Logger
}

func NewNotesController(store NotesStore, logger Logger) NotesCtrlr {
	return NotesCtrlr{NotesStore: store, Logger: logger}
}

func (ns *NotesCtrlr) ProcessAddNote(w http.ResponseWriter, r *http.Request) {
	ns.Logger.Infof("%s request to %s received", r.Method, r.URL.Path)
	var body map[string]Note
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ns.Logger.Errorf("Error Unmarshaling request body")
		return
	}

	note, ok := body["note"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		ns.Logger.Errorf("Invalid body")
		return
	}

	_ = ns.NotesStore.AddNote(note)
	w.WriteHeader(http.StatusAccepted)
}

func (ns *NotesCtrlr) GetNotesByID(w http.ResponseWriter, r *http.Request) {
	ns.Logger.Infof("%s request to %s received", r.Method, r.URL.Path)
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
	ns.Logger.Infof("%s request to %s received", r.Method, r.URL.Path)
	notes := ns.NotesStore.GetAllNotes()
	json.NewEncoder(w).Encode(notes)
	return
}
