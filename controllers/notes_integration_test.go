package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddingNotesAndRetrievingThem(t *testing.T) {
	store := NewInMemoryNotesStore()
	notesC := Notes{&store}
	userID := 1

	notesC.ProcessAddNote(httptest.NewRecorder(), newPostAddNoteRequest(userID, "Test note 1"))
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostAddNoteRequest(userID, "Test note 2"))
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostAddNoteRequest(userID, "Test note 3"))

	response := httptest.NewRecorder()
	notesC.GetNotesByID(response, newGetNotesByUserIdRequest(userID))
	assertStatusCode(t, response.Result().StatusCode, http.StatusOK)
	assertResponseBody(t, response.Body, []string{"Test note 1", "Test note 2", "Test note 3"})
}

type InMemoryNotesStore struct {
	notes map[int][]string
}

func NewInMemoryNotesStore() InMemoryNotesStore {
	notes := make(map[int][]string)
	return InMemoryNotesStore{notes: notes}
}

func (i *InMemoryNotesStore) GetNotesByID(id int) []string {
	return i.notes[id]
}

func (i *InMemoryNotesStore) GetAllNotes() []string {
	var allNotes []string
	for _, notes := range i.notes {
		allNotes = append(allNotes, notes...)
	}
	return allNotes
}

func (i *InMemoryNotesStore) AddNote(userID int, note string) error {
	// if _, ok := i.notes[userID]; !ok {
	// 	i.notes[userID] = []string{note}
	// 	return nil
	// }
	i.notes[userID] = append(i.notes[userID], note)
	return nil
}
