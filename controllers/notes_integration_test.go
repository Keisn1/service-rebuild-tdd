package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddingNotesAndRetrievingThem(t *testing.T) {
	store := NewInMemoryNotesStore()
	logger := NewSimpleLogger()

	notesC := NotesCtrlr{NotesStore: &store, Logger: &logger}

	userID := 1
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostRequestWithNote(t, NewNote(userID, "Test note 1"), "/notes/1"))
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostRequestWithNote(t, NewNote(userID, "Test note 2"), "/notes/1"))
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostRequestWithNote(t, NewNote(userID, "Test note 3"), "/notes/1"))

	userID = 2
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostRequestWithNote(t, NewNote(userID, "Test note 4"), "/notes/2"))
	notesC.ProcessAddNote(httptest.NewRecorder(), newPostRequestWithNote(t, NewNote(userID, "Test note 5"), "/notes/2"))

	// Testing notes by id
	assertNotesByIdAsExpected(t, 1, Notes{{1, "Test note 1"}, {1, "Test note 2"}, {1, "Test note 3"}}, notesC)
	assertNotesByIdAsExpected(t, 2, Notes{{2, "Test note 4"}, {2, "Test note 5"}}, notesC)

	// Testing all notes
	wantAllNotes := Notes{{1, "Test note 1"}, {1, "Test note 2"}, {1, "Test note 3"}, {2, "Test note 4"}, {2, "Test note 5"}}
	assertAllNotesAsExpected(t, wantAllNotes, notesC)
}

func assertNotesByIdAsExpected(t testing.TB, userID int, wantNotes Notes, notesC NotesCtrlr) {
	t.Helper()
	response := httptest.NewRecorder()
	notesC.GetNotesByUserID(response, newGetNotesByUserIdRequest(t, userID))

	gotNotes := getNotesFromResponse(t, response.Body)
	assertStatusCode(t, response.Result().StatusCode, http.StatusOK)
	assertNotesEqual(t, gotNotes, wantNotes)
}

func assertAllNotesAsExpected(t testing.TB, wantAllNotes Notes, notesC NotesCtrlr) {
	t.Helper()
	response := httptest.NewRecorder()
	notesC.GetAllNotes(response, newGetAllNotesRequest(t))

	gotAllNotes := getNotesFromResponse(t, response.Body)
	assertNotesEqual(t, gotAllNotes, wantAllNotes)
}

func getNotesFromResponse(t testing.TB, body io.Reader) (notes Notes) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&notes)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into map[int]Notes", err)
	}
	return
}
