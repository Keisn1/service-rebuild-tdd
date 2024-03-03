package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"errors"

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
	GetNotesByUserID(int) Notes
	AddNote(Note) error
	Delete(int) error
}

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, a ...any)
}

type NotesCtrlr struct {
	NotesStore NotesStore
	Logger     Logger
}

func NewNotesController(store NotesStore, logger Logger) NotesCtrlr {
	return NotesCtrlr{NotesStore: store, Logger: logger}
}

var (
	ErrUnmarshalRequestBody = errors.New("Error Unmarshaling request body")
	ErrDBResourceCreation   = errors.New("Could not create resource")
	ErrDBResourceDeletion   = errors.New("Could not delete resource")
	ErrInvalidRequestBody   = errors.New("Invalid request body")
	ErrInvalidUserID        = errors.New("Invalid user ID")
	ErrInvalidNoteId        = errors.New("Invalid note ID")
)

func (ns *NotesCtrlr) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ns.Logger.Errorf("%w: %w", ErrInvalidNoteId, err)
		return
	}

	fmt.Println(id)
	err = ns.NotesStore.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		ns.Logger.Errorf("%w: %w", ErrDBResourceDeletion, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ns *NotesCtrlr) ProcessAddNote(w http.ResponseWriter, r *http.Request) {
	ns.Logger.Infof("%s request to %s received", r.Method, r.URL.Path)
	var body map[string]Note
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ns.Logger.Errorf("%w: %w", ErrUnmarshalRequestBody, err)
		return
	}

	note, ok := body["note"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		ns.Logger.Errorf("%w: %w", ErrInvalidRequestBody, err)
		return
	}

	err = ns.NotesStore.AddNote(note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ns.Logger.Errorf("%w: %w", ErrDBResourceCreation, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (ns *NotesCtrlr) GetNotesByUserID(w http.ResponseWriter, r *http.Request) {
	ns.Logger.Infof("%s request to %s received", r.Method, r.URL.Path)

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ns.Logger.Errorf("%w: %w", ErrInvalidUserID, err)
		return
	}

	notes := ns.NotesStore.GetNotesByUserID(userID)
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
