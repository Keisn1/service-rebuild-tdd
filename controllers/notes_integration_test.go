package controllers_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strconv"
// 	"testing"

// 	"fmt"
// 	ctrls "github.com/Keisn1/note-taking-app/controllers"
// 	"github.com/Keisn1/note-taking-app/domain"
// 	"github.com/stretchr/testify/assert"
// )

// func buildNotes(t testing.TB, userID int, notes []string) domain.Notes {
// 	t.Helper()
// 	var ret domain.Notes
// 	for _, note := range notes {
// 		ret = append(ret, domain.Note{UserID: userID, Note: note})
// 	}
// 	return ret
// }

// func allNotes() domain.Notes {
// 	return domain.Notes{
// 		{UserID: 1, Note: "Test note 1"},
// 		{UserID: 1, Note: "Test note 2"},
// 		{UserID: 1, Note: "Test note 3"},
// 		{UserID: 2, Note: "Test note 4"},
// 		{UserID: 2, Note: "Test note 5"},
// 		{UserID: 2, Note: "Test note 6"},
// 	}
// }

// func TestIntegration(t *testing.T) {
// 	store := ctrls.NewInMemoryNotesStore()
// 	notesC := ctrls.NewNotesCtrlr(store)

// 	// Add notes
// 	addNotes(t, notesC)

// 	// Testing all notes
// 	canRetrieveAllNotes(t, notesC)

// 	// Testing notes by userID
// 	canRetrieveNotesByUserID(t, notesC)

// 	// Testing notes by userID and noteID
// 	canRetrieveNotesByUserIDAndNoteID(t, notesC)

// 	// Edit a note
// 	canEditNote(t, notesC)

// 	// Delete a note
// 	canDeleteNotes(t, notesC)
// }

// func canDeleteNotes(t *testing.T, notesC ctrls.NotesCtrlr) {
// 	t.Helper()

// 	rr := httptest.NewRecorder()
// 	req := setupRequest(t, "GET", "/users/notes", urlParams{}, &bytes.Buffer{})
// 	notesC.GetAllNotes(rr, req)
// 	allNotes := decodeBodyNotes(t, rr.Body)

// 	for _, n := range allNotes {
// 		up := urlParams{userID: strconv.Itoa(n.UserID), noteID: strconv.Itoa(n.NoteID)}
// 		rr := httptest.NewRecorder()
// 		req := setupRequest(t, "DELETE", "/users/{userID}/notes/{noteID}", up, &bytes.Buffer{})
// 		notesC.Delete(rr, req)
// 		assert.Equal(t, http.StatusNoContent, rr.Code)

// 		rr = httptest.NewRecorder()
// 		req = setupRequest(t, "GET", "/users/{userID}/notes/{noteID}", up, &bytes.Buffer{})
// 		notesC.GetNoteByUserIDAndNoteID(rr, req)

// 		// gotNotes := decodeBodyNotes(t, rr.Body)
// 		assert.Equal(t, http.StatusNotFound, rr.Code)
// 	}

// 	rr = httptest.NewRecorder()
// 	req = setupRequest(t, "GET", "/users/notes", urlParams{}, &bytes.Buffer{})
// 	notesC.GetAllNotes(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	allNotes = decodeBodyNotes(t, rr.Body)
// 	assert.Equal(t, 0, len(allNotes))
// }

// func canEditNote(t *testing.T, notesC ctrls.NotesCtrlr) {
// 	rr := httptest.NewRecorder()
// 	req := setupRequest(t, "GET", "/users/notes", urlParams{}, &bytes.Buffer{})
// 	notesC.GetAllNotes(rr, req)
// 	allNotes := decodeBodyNotes(t, rr.Body)

// 	for _, n := range allNotes {
// 		n.Note = fmt.Sprintf("Edited note userID %v noteID %v", n.UserID, n.NoteID)

// 		// edit note
// 		up := urlParams{userID: strconv.Itoa(n.UserID), noteID: strconv.Itoa(n.NoteID)}
// 		body := domain.NotePost{Note: n.Note}
// 		rr := httptest.NewRecorder()
// 		req := setupRequest(t, "POST", "/users/{userID}/notes/{noteID}", up, mustEncode(t, body))
// 		notesC.Edit(rr, req)
// 		assert.Equal(t, http.StatusAccepted, rr.Code)

// 		// test
// 		rr = httptest.NewRecorder()
// 		req = setupRequest(t, "GET", "/users/{userID}/notes/{noteID}", up, &bytes.Buffer{})
// 		notesC.GetNoteByUserIDAndNoteID(rr, req)
// 		gotNotes := decodeBodyNotes(t, rr.Body)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		assert.Equal(t, domain.Notes{n}, gotNotes)
// 	}
// }

// func canRetrieveNotesByUserIDAndNoteID(t *testing.T, notesC ctrls.NotesCtrlr) {
// 	t.Helper()

// 	rr := httptest.NewRecorder()
// 	req := setupRequest(t, "GET", "/users/notes", urlParams{}, &bytes.Buffer{})
// 	notesC.GetAllNotes(rr, req)
// 	allNotes := decodeBodyNotes(t, rr.Body)

// 	for _, n := range allNotes {
// 		rr := httptest.NewRecorder()
// 		req := setupRequest(t, "POST", "/users/{userID}/notes/{noteID}", urlParams{
// 			userID: strconv.Itoa(n.UserID),
// 			noteID: strconv.Itoa(n.NoteID),
// 		}, &bytes.Buffer{})
// 		notesC.GetNoteByUserIDAndNoteID(rr, req)
// 		gotNotes := decodeBodyNotes(t, rr.Body)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		assert.Equal(t, domain.Notes{n}, gotNotes)
// 	}
// }

// func canRetrieveNotesByUserID(t *testing.T, notesC ctrls.NotesCtrlr) {
// 	t.Helper()
// 	for i := range []int{1, 2} {
// 		var wantNotes domain.Notes
// 		for _, n := range allNotes() {
// 			if n.UserID == i {
// 				wantNotes = append(wantNotes, n)
// 			}
// 		}

// 		rr := httptest.NewRecorder()
// 		req := setupRequest(t, "GET", "/users/notes", urlParams{userID: strconv.Itoa(i)}, &bytes.Buffer{})
// 		notesC.GetNotesByUserID(rr, req)
// 		gotNotes := decodeBodyNotes(t, rr.Body)
// 		gotNotes = setNoteIdZero(t, gotNotes)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		assert.Equal(t, wantNotes, gotNotes)
// 	}
// }

// func canRetrieveAllNotes(t *testing.T, notesC ctrls.NotesCtrlr) {
// 	t.Helper()
// 	rr := httptest.NewRecorder()
// 	req := setupRequest(t, "GET", "/users/notes", urlParams{}, &bytes.Buffer{})
// 	notesC.GetAllNotes(rr, req)

// 	gotNotes := decodeBodyNotes(t, rr.Body)
// 	gotNotes = setNoteIdZero(t, gotNotes)
// 	assert.Equal(t, allNotes(), gotNotes)
// }

// func setNoteIdZero(t *testing.T, notes domain.Notes) domain.Notes {
// 	var newNotes domain.Notes
// 	for _, n := range notes {
// 		n.NoteID = 0
// 		newNotes = append(newNotes, n)
// 	}
// 	return newNotes
// }

// func decodeBodyNotes(t testing.TB, body io.Reader) (notes domain.Notes) {
// 	t.Helper()
// 	err := json.NewDecoder(body).Decode(&notes)
// 	if err != nil {
// 		t.Fatalf("Unable to parse body into Notes: %v", err)
// 	}
// 	return
// }

// func addNotes(t *testing.T, notesC ctrls.NotesCtrlr) {
// 	t.Helper()
// 	for _, n := range allNotes() {
// 		body := domain.NotePost{Note: n.Note}
// 		req := setupRequest(t, "POST", "/users/{userID}/notes", urlParams{userID: strconv.Itoa(n.UserID)}, mustEncode(t, body))
// 		notesC.Add(
// 			httptest.NewRecorder(),
// 			req,
// 		)
// 	}
// }
