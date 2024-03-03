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
	AddNote(userID int, note string) error
	EditNote(userID, noteID int, note string) error
	Delete(userID int, noteID int) error
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
	ErrEncodingJson         = errors.New("Error encoding json")
	ErrDBResourceCreation   = errors.New("Could not create resource")
	ErrDBResourceDeletion   = errors.New("Could not delete resource")
	ErrInvalidRequestBody   = errors.New("Invalid request body")
	ErrInvalidUserID        = errors.New("Invalid user ID")
	ErrInvalidNoteID        = errors.New("Invalid note ID")
)

func (nc *NotesCtrlr) Edit(w http.ResponseWriter, r *http.Request) {
	// var body map[string]Note
	// _ = json.NewDecoder(r.Body).Decode(&body)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	nc.Logger.Errorf("%w: %w", ErrUnmarshalRequestBody, err)
	// 	return
	// }

	// note, _ := body["note"]
	// if !ok {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	nc.Logger.Errorf("%w: %w", ErrInvalidRequestBody, err)
	// 	return
	// }

	// _ = nc.NotesStore.EditNote(userID, noteID, note)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	nc.Logger.Errorf("%w: %w", ErrDBResourceCreation, err)
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
}

func (nc *NotesCtrlr) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("%w: %w", ErrInvalidUserID, err)
		return
	}

	noteID, err := strconv.Atoi(chi.URLParam(r, "noteID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("%w: %w", ErrInvalidNoteID, err)
		return
	}

	err = nc.NotesStore.Delete(userID, noteID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		nc.Logger.Errorf("Delete DBError: %w", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	nc.Logger.Errorf("Success: Delete noteID %v userID %v")
}

func (nc *NotesCtrlr) ProcessAddNote(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("ProcessAddNote invalid userID: %w", err)
		return
	}

	var body map[string]string
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("ProcessAddNote invalid json: %w", err)
		return
	}

	note, ok := body["note"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("ProcessAddNote invalid body: %w", err)
		return
	}

	err = nc.NotesStore.AddNote(userID, note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		nc.Logger.Errorf("ProcessAddNote DBerror: %w", err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	nc.Logger.Infof("Success: ProcessAddNote with userID %d and note %v", userID, note)
}

func (nc *NotesCtrlr) GetNotesByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("GetNotesByUserID invalid userID: %w", err)
		return
	}

	notes := nc.NotesStore.GetNotesByUserID(userID)
	if len(notes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		nc.Logger.Errorf("GetNotesByUserID user not Found: %w", userID)
		return
	}

	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		nc.Logger.Errorf("GetAllNotes invalid json: %w", err)
		return
	}

	nc.Logger.Infof("Success: GetNotesByUserID with userID %d", userID)
	return
}

func (nc *NotesCtrlr) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes := nc.NotesStore.GetAllNotes()
	err := json.NewEncoder(w).Encode(notes)
	if err != nil {
		nc.Logger.Errorf("GetAllNotes invalid json: %w", err)
		return
	}
	nc.Logger.Infof("Success: GetAllNotes")
	return
}
