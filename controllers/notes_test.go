package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	"context"

	"github.com/go-chi/chi"
)

type StubNotesStore struct {
	notes        map[int][]string
	addNoteCalls []addNoteCall
}

type addNoteCall struct {
	userID int
	note   string
}

func (sns *StubNotesStore) AddNote(userID int, note string) error {
	sns.addNoteCalls = append(sns.addNoteCalls, addNoteCall{userID, note})
	return nil
}

func (sns *StubNotesStore) GetAllNotes() []string {
	var allNotes []string
	for _, notes := range sns.notes {
		allNotes = append(allNotes, notes...)
	}
	return allNotes
}

func (sns *StubNotesStore) GetNotesByID(id int) []string {
	return sns.notes[id]
}

func TestNotes(t *testing.T) {
	notesStore := StubNotesStore{
		notes: map[int][]string{
			1: {"Note 1 user 1", "Note 2 user 1"},
			2: {"Note 1 user 2", "Note 2 user 2"},
		},
	}
	notesC := &Notes{NotesStore: &notesStore}

	t.Run("Server returns all Notes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/notes", nil)
		response := httptest.NewRecorder()
		notesC.GetAllNotes(response, request)

		assertStatusCode(t, response.Result().StatusCode, http.StatusOK)

		got := decodeJsonBody(response.Body)
		want := []string{"Note 1 user 1", "Note 2 user 1", "Note 1 user 2", "Note 2 user 2"}
		assertStringSlicesAreEqual(t, got, want)
	})

	t.Run("Return notes for user with userID", func(t *testing.T) {
		testCases := []struct {
			userID     int
			want       []string
			statusCode int
		}{
			{1, []string{"Note 1 user 1", "Note 2 user 1"}, http.StatusOK},
			{2, []string{"Note 1 user 2", "Note 2 user 2"}, http.StatusOK},
			{100, []string{}, http.StatusNotFound},
		}

		for _, tc := range testCases {
			response := httptest.NewRecorder()
			request := newGetNotesByUserIdRequest(tc.userID)
			notesC.GetNotesByID(response, request)

			assertStatusCode(t, response.Result().StatusCode, tc.statusCode)
			assertResponseBody(t, response.Body, tc.want)
		}
	})

	t.Run("adds a note with POST", func(t *testing.T) {
		userID := 1
		testNote := "Test note"
		request := newPostAddNoteRequest(userID, testNote)
		response := httptest.NewRecorder()
		notesC.ProcessAddNote(response, request)

		// assertions
		assertStatusCode(t, response.Result().StatusCode, http.StatusAccepted)
		assertLengthSlice(t, notesStore.addNoteCalls, 1)

		want := addNoteCall{userID, testNote}
		got := notesStore.addNoteCalls[0]
		if !reflect.DeepEqual(got, want) {
			t.Errorf(`got = %v; want %v`, got, want)
		}
	})
}

func newPostAddNoteRequest(userID int, note string) *http.Request {
	url := fmt.Sprintf("/notes/%d", userID)
	requestBody := map[string]string{"note": note}
	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(requestBody)
	request, _ := http.NewRequest(http.MethodPost, url, buf)
	return WithUrlParam(request, "id", fmt.Sprintf("%v", userID))
}

func newGetNotesByUserIdRequest(userID int) *http.Request {
	request, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/notes/%v", userID),
		nil)
	return WithUrlParam(request, "id", fmt.Sprintf("%v", userID))
}

func assertLengthSlice[T any](t testing.TB, elements []T, want int) {
	t.Helper()
	if len(elements) != want {
		t.Errorf(`got = %v; want %v`, len(elements), want)
	}

}

func assertStringSlicesAreEqual(t testing.TB, got, want []string) {
	t.Helper()
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	sort.Slice(want, func(i, j int) bool { return want[i] < want[j] })
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func decodeJsonBody(r io.Reader) []string {
	var res []string
	json.NewDecoder(r).Decode(&res)
	return res
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func assertSlicesHaveSameLength[T any](t testing.TB, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf(`len(got) = %v; len(want) %v`, len(got), len(want))
	}
}

// WithUrlParam returns a pointer to a request object with the given URL params
// added to a new chi.Context object.
func WithUrlParam(r *http.Request, key, value string) *http.Request {
	chiCtx := chi.NewRouteContext()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add(key, value)
	return r
}

func assertResponseBody(t testing.TB, body *bytes.Buffer, want []string) {
	got := decodeJsonBody(body)
	assertSlicesHaveSameLength(t, got, want)
	assertStringSlicesAreEqual(t, got, want)
}
