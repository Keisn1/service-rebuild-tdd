package controllers

import (
	"encoding/json"
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
	EditNote(Note) error
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

func NewNotesCtrlr(store NotesStore, logger Logger) NotesCtrlr {
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

func (nc *NotesCtrlr) Edit(w http.ResponseWriter, r *http.Request) {
	nc.Logger.Infof("%s request to %s received", r.Method, r.URL.Path)
	var body map[string]Note
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("%w: %w", ErrUnmarshalRequestBody, err)
		return
	}

	note, ok := body["note"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("%w: %w", ErrInvalidRequestBody, err)
		return
	}

	err = nc.NotesStore.EditNote(note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		nc.Logger.Errorf("%w: %w", ErrDBResourceCreation, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (nc *NotesCtrlr) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("%w: %w", ErrInvalidNoteId, err)
		return
	}

	err = nc.NotesStore.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		nc.Logger.Errorf("%w: %w", ErrDBResourceDeletion, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (nc *NotesCtrlr) ProcessAddNote(w http.ResponseWriter, r *http.Request) {
	nc.Logger.Infof("%s request to %s received", r.Method, r.URL.Path)
	var body map[string]Note
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("%w: %w", ErrUnmarshalRequestBody, err)
		return
	}

	note, ok := body["note"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("%w: %w", ErrInvalidRequestBody, err)
		return
	}

	err = nc.NotesStore.AddNote(note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		nc.Logger.Errorf("%w: %w", ErrDBResourceCreation, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (nc *NotesCtrlr) GetNotesByUserID(w http.ResponseWriter, r *http.Request) {
	nc.Logger.Infof("%s request to %s received", r.Method, r.URL.Path)

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("%w: %w", ErrInvalidUserID, err)
		return
	}

	notes := nc.NotesStore.GetNotesByUserID(userID)
	if len(notes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode([]string{})
		return
	}
	json.NewEncoder(w).Encode(notes)
	return
}

func (nc *NotesCtrlr) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	nc.Logger.Infof("%s request to %s received", r.Method, r.URL.Path)
	notes := nc.NotesStore.GetAllNotes()
	json.NewEncoder(w).Encode(notes)
	return
}
