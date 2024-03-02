package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddingNotesAndRetrievingThem(t *testing.T) {
	store := NewInMemoryNotesStore()
	logger := NewSimpleLogger()

	notesC := NotesCtrlr{NotesStore: &store, Logger: &logger}

	userID := 1
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostRequestWithNote(t, NewNote(userID, "Test note 1")))
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostRequestWithNote(t, NewNote(userID, "Test note 2")))
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostRequestWithNote(t, NewNote(userID, "Test note 3")))

	userID = 2
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostRequestWithNote(t, NewNote(userID, "Test note 4")))
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostRequestWithNote(t, NewNote(userID, "Test note 5")))

	// Testing notes by id
	assertNotesByIdAsExpected(t, 1, Notes{{1, "Test note 1"}, {1, "Test note 2"}, {1, "Test note 3"}}, notesC)
	assertNotesByIdAsExpected(t, 2, Notes{{2, "Test note 4"}, {2, "Test note 5"}}, notesC)

	wantAllNotes := Notes{{1, "Test note 1"}, {1, "Test note 2"}, {1, "Test note 3"}, {2, "Test note 4"}, {2, "Test note 5"}}
	assertAllNotesAsExpected(t, wantAllNotes, notesC)
}

func assertNotesByIdAsExpected(t testing.TB, userID int, wantNotes Notes, notesC NotesCtrlr) {
	t.Helper()
	response := httptest.NewRecorder()
	notesC.GetNotesByID(response, newGetNotesByUserIdRequest(t, userID))

	gotNotes := getNotesByIdFromResponse(t, response.Body)
	assertStatusCode(t, response.Result().StatusCode, http.StatusOK)
	assertNotesById(t, gotNotes, wantNotes)
}

func assertAllNotesAsExpected(t testing.TB, wantAllNotes Notes, notesC NotesCtrlr) {
	t.Helper()
	response := httptest.NewRecorder()
	notesC.GetAllNotes(response, newGetAllNotesRequest(t))

	gotAllNotes := getAllNotesFromResponse(t, response.Body)
	assertAllNotes(t, gotAllNotes, wantAllNotes)
}
