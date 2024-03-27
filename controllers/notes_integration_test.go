package controllers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Keisn1/note-taking-app/common"
	ctrls "github.com/Keisn1/note-taking-app/controllers"
	"github.com/Keisn1/note-taking-app/domain"
)

func buildNotes(t testing.TB, userID int, notes []string) domain.Notes {
	t.Helper()
	var ret domain.Notes
	for _, note := range notes {
		ret = append(ret, domain.Note{UserID: userID, Note: note})
	}
	return ret
}

func TestIntegrationInMemoryStore(t *testing.T) {
	store := ctrls.NewInMemoryNotesStore()
	logger := common.NewSimpleLogger()
	notesC := ctrls.NewNotesCtrlr(store, &logger)

	notesUser1 := buildNotes(t, 1, []string{"Test note 1", "Test note 2", "Test note 3"})
	notesUser2 := buildNotes(t, 2, []string{"Test note 4", "Test note 5"})
	allNotes := append(notesUser1, notesUser2...)

	addNotes(t, notesUser1, notesC)
	addNotes(t, notesUser2, notesC)

	// Testing all notes
	assertAllNotesAsExpected(t, allNotes, notesC)

	// Testing notes by userID
	assertNotesByIdAsExpected(t, 1, notesUser1, notesC)
	assertNotesByIdAsExpected(t, 2, notesUser2, notesC)

	// Edit a note
	pNote, text := EditNote(t, notesC)
	assertNoteWasEdited(t, pNote, text, notesC)

	// Delete a note
	restOfNotes := deleteANote(t, notesC)
	assertAllNotesAsExpected(t, restOfNotes, notesC)
}

func deleteANote(t testing.TB, notesC ctrls.NotesCtrlr) (restOfNotes domain.Notes) {
	t.Helper()
	response := httptest.NewRecorder()
	notesC.GetAllNotes(response, newGetAllNotesRequest(t))
	allNotes := decodeBodyNotes(t, response.Body)
	dNote, restOfNotes := allNotes[0], allNotes[1:]

	deleteRequest, err := http.NewRequest(http.MethodDelete, "", nil)
	assertNoError(t, err)
	deleteRequest = WithUrlParams(deleteRequest, Params{
		"userID": strconv.Itoa(dNote.UserID),
		"noteID": strconv.Itoa(dNote.NoteID),
	})
	notesC.Delete(response, deleteRequest)
	return restOfNotes
}

func assertNoteWasEdited(t testing.TB, pNote domain.Note, text string, notesC ctrls.NotesCtrlr) {
	response := httptest.NewRecorder()
	notesC.GetAllNotes(response, newGetAllNotesRequest(t))
	allNotes := decodeBodyNotes(t, response.Body)
	for _, n := range allNotes {
		if n.UserID == pNote.UserID && n.NoteID == pNote.NoteID && n.Note != text {
			t.Errorf("Did not edit note with userID %d, noteID %d and note %s to \"Edit Note\"", n.UserID, n.NoteID, n.Note)
		}
	}
}

func assertNotesByIdAsExpected(t testing.TB, userID int, wantNotes domain.Notes, notesC ctrls.NotesCtrlr) {
	t.Helper()
	response := httptest.NewRecorder()
	notesC.GetNotesByUserID(response, newGetNotesByUserIdRequest(t, userID))

	gotNotes := decodeBodyNotes(t, response.Body)
	assertStatusCode(t, response.Result().StatusCode, http.StatusOK)
	compareNotesByUserIDAndNote(t, gotNotes, wantNotes)
}

func EditNote(t testing.TB, notesC ctrls.NotesCtrlr) (domain.Note, string) {
	// returns note that is being edited
	response := httptest.NewRecorder()
	notesC.GetAllNotes(response, newGetAllNotesRequest(t))
	allNotes := decodeBodyNotes(t, response.Body)

	note := allNotes[0]
	notesC.Edit(httptest.NewRecorder(), newPutRequestWithNoteAndUrlParams(t, "Edit Note", Params{
		"userID": strconv.Itoa(note.UserID),
		"noteID": strconv.Itoa(note.NoteID),
	}))
	return note, "Edit Note"
}

// compareNotesByUserIDAndNote compares two slices of Notes by UserID and Note fields
func compareNotesByUserIDAndNote(t testing.TB, got, want domain.Notes) {
	t.Helper()
	assertSlicesSameLength(t, got, want)

	// Create a map for the wantNotes for quick lookups
	for _, w := range want {
		found := false
		for _, g := range got {
			if w.UserID == g.UserID && w.Note == g.Note {
				found = true
			}
		}
		if !found {
			t.Errorf("got %v not equal want %v", got, want)
		}
	}
}

func assertAllNotesAsExpected(t testing.TB, wantNotes domain.Notes, notesC ctrls.NotesCtrlr) {
	t.Helper()
	response := httptest.NewRecorder()
	notesC.GetAllNotes(response, newGetAllNotesRequest(t))

	gotNotes := decodeBodyNotes(t, response.Body)
	compareNotesByUserIDAndNote(t, gotNotes, wantNotes)
}

func decodeBodyNotes(t testing.TB, body io.Reader) (notes domain.Notes) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&notes)
	if err != nil {
		t.Fatalf("Unable to parse body into Notes: %v", err)
	}
	return
}

func addNotes(t testing.TB, notes domain.Notes, notesC ctrls.NotesCtrlr) {
	for _, n := range notes {
		notesC.Add(
			httptest.NewRecorder(),
			newPostRequestWithNoteAndUrlParam(t, n.Note, "userID", strconv.Itoa(n.UserID)),
		)
	}
}
