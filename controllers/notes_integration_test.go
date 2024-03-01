package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddingNotesAndRetrievingThem(t *testing.T) {
	data := make(map[int][]string)
	store := NewInMemoryNotesStore(data)
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
