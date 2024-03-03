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
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if handleBadRequest(w, err, nc.Logger, "Edit", "userID") {
		return
	}

	noteID, err := strconv.Atoi(chi.URLParam(r, "noteID"))
	if handleBadRequest(w, err, nc.Logger, "Edit", "noteID") {
		return
	}

	var body map[string]string
	err = json.NewDecoder(r.Body).Decode(&body)
	if handleBadRequest(w, err, nc.Logger, "Edit", "json") {
		return
	}

	note, ok := body["note"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("ProcessAddNote invalid body: %v", err)
		return
	}

	err = nc.NotesStore.EditNote(userID, noteID, note)
	if handleError(w, err, http.StatusInternalServerError, "Edit", "DBerror", nc.Logger) {
		return
	}
	w.WriteHeader(http.StatusOK)
	nc.Logger.Infof("Success: Edit: userID %v noteID %v note %v", userID, noteID, note)
}

func (nc *NotesCtrlr) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if handleBadRequest(w, err, nc.Logger, "Delete", "userID") {
		return
	}

	noteID, err := strconv.Atoi(chi.URLParam(r, "noteID"))
	if handleBadRequest(w, err, nc.Logger, "Delete", "noteID") {
		return
	}

	err = nc.NotesStore.Delete(userID, noteID)
	if handleError(w, err, http.StatusInternalServerError, "Delete", "DBerror", nc.Logger) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
	nc.Logger.Infof("Success: Delete noteID %v userID %v", noteID, userID)
}

func (nc *NotesCtrlr) ProcessAddNote(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if handleBadRequest(w, err, nc.Logger, "ProcessAddNote", "userID") {
		return
	}

	var body map[string]string
	err = json.NewDecoder(r.Body).Decode(&body)
	if handleBadRequest(w, err, nc.Logger, "ProcessAddNote", "json") {
		return
	}

	note, ok := body["note"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		nc.Logger.Errorf("ProcessAddNote invalid body: %w", err)
		return
	}

	err = nc.NotesStore.AddNote(userID, note)
	if handleError(w, err, http.StatusConflict, "ProcessAddNote", "DBerror", nc.Logger) {
		return
	}

	w.WriteHeader(http.StatusAccepted)
	nc.Logger.Infof("Success: ProcessAddNote with userID %d and note %v", userID, note)
}

func (nc *NotesCtrlr) GetNotesByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if handleBadRequest(w, err, nc.Logger, "GetNotesByUserID", "userID") {
		return
	}

	notes := nc.NotesStore.GetNotesByUserID(userID)
	if len(notes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		nc.Logger.Errorf("GetNotesByUserID user not Found: %w", userID)
		return
	}

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		nc.Logger.Errorf("GetAllNotes invalid json: %w", err)
		return
	}

	nc.Logger.Infof("Success: GetNotesByUserID with userID %d", userID)
	return
}

func (nc *NotesCtrlr) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes := nc.NotesStore.GetAllNotes()
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		nc.Logger.Errorf("GetAllNotes invalid json: %w", err)
		return
	}
	nc.Logger.Infof("Success: GetAllNotes")
	return
}

func handleBadRequest(w http.ResponseWriter, err error, logger Logger, action, param string) bool {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Errorf("%s invalid %s: %v", action, param, err)
		return true
	}
	return false
}

func handleError(w http.ResponseWriter, err error, httpErr uint, action, msg string, logger Logger) bool {
	if err != nil {
		w.WriteHeader(int(httpErr))
		logger.Errorf("%s %s: %w", action, msg, err)
		return true
	}
	return false
}
