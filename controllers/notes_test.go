package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	"github.com/go-chi/chi"
)

type StubNotesStore struct {
	notes        Notes
	addNoteCalls Notes
}

func (sns *StubNotesStore) AddNote(note Note) error {
	sns.addNoteCalls = append(sns.addNoteCalls, note)
	return nil
}

func (sns *StubNotesStore) GetAllNotes() Notes {
	return sns.notes
}

func (sns *StubNotesStore) GetNotesByID(userID int) (ret Notes) {
	for _, n := range sns.notes {
		if n.UserID == userID {
			ret = append(ret, n)
		}
	}
	return
}

type StubLogger struct {
	infolnCalls []any
}

func (sl *StubLogger) Infoln(a ...any) {
	for _, val := range a {
		sl.infolnCalls = append(sl.infolnCalls, val)
	}
}

func (sl *StubLogger) Reset() {
	sl.infolnCalls = []any{}
}

func TestNotes(t *testing.T) {
	notesStore := StubNotesStore{
		notes: Notes{
			{1, "Note 1 user 1"}, {1, "Note 2 user 1"},
			{2, "Note 1 user 2"}, {2, "Note 2 user 2"},
		},
	}
	logger := StubLogger{}
	notesC := &NotesCtrlr{NotesStore: &notesStore, Logger: &logger}

	t.Run("Server returns all Notes", func(t *testing.T) {
		logger.Reset()
		wantedNotes := Notes{
			{1, "Note 1 user 1"}, {1, "Note 2 user 1"},
			{2, "Note 1 user 2"}, {2, "Note 2 user 2"},
		}

		request := newGetAllNotesRequest(t)
		response := httptest.NewRecorder()
		notesC.GetAllNotes(response, request)

		got := getAllNotesFromResponse(t, response.Body)
		assertStatusCode(t, response.Result().StatusCode, http.StatusOK)
		assertAllNotes(t, got, wantedNotes)
		assertLoggerInfolnCalls(t, logger.infolnCalls, []string{"GET request to /notes received"})
	})

	t.Run("Return notes for user with userID", func(t *testing.T) {
		logger.Reset()
		testCases := []struct {
			userID     int
			want       Notes
			statusCode int
		}{
			{1, Notes{{1, "Note 1 user 1"}, {1, "Note 2 user 1"}}, http.StatusOK},
			{2, Notes{{2, "Note 1 user 2"}, {2, "Note 2 user 2"}}, http.StatusOK},
			{100, Notes{}, http.StatusNotFound},
		}

		for _, tc := range testCases {
			response := httptest.NewRecorder()
			request := newGetNotesByUserIdRequest(t, tc.userID)
			notesC.GetNotesByID(response, request)

			assertStatusCode(t, response.Result().StatusCode, tc.statusCode)

			got := getNotesByIdFromResponse(t, response.Body)
			assertNotesById(t, got, tc.want)
		}
		assertLoggerInfolnCalls(t, logger.infolnCalls, []string{
			"GET request to /notes/1 received",
			"GET request to /notes/2 received",
			"GET request to /notes/100 received",
		})
	})

	t.Run("adds a note with POST", func(t *testing.T) {
		logger.Reset()
		note := NewNote(1, "Test note")
		wantAddNoteCalls := Notes{note}
		request := newPostAddNoteRequest(t, note)
		response := httptest.NewRecorder()
		notesC.ProcessAddNote(response, request)

		// assertions
		assertStatusCode(t, response.Result().StatusCode, http.StatusAccepted)
		assertAddNoteCalls(t, notesStore.addNoteCalls, wantAddNoteCalls)
		assertLoggerInfolnCalls(t, logger.infolnCalls, []string{"POST request to /notes/1 received"})
	})
}

func newPostAddNoteRequest(t testing.TB, note Note) *http.Request {
	t.Helper()
	url := fmt.Sprintf("/notes/%d", note.UserID)
	requestBody := map[string]Note{"note": note}
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(requestBody)
	if err != nil {
		t.Fatalf("Error encoding requestBody when build addNote post request %q", err)
	}
	request, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		t.Fatalf("Could not build request newPostAddNoteRequest: %q", err)
	}
	return WithUrlParam(request, "id", fmt.Sprintf("%v", note.UserID))
}

func newGetNotesByUserIdRequest(t testing.TB, userID int) *http.Request {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/notes/%v", userID),
		nil)
	if err != nil {
		t.Fatalf("Could not build request newPostAddNoteRequest: %q", err)
	}
	return WithUrlParam(request, "id", fmt.Sprintf("%v", userID))
}

func assertLengthSlice[T any](t testing.TB, elements []T, want int) {
	t.Helper()
	if len(elements) != want {
		t.Errorf(`got = %v; want %v`, len(elements), want)
	}

}

func assertStringSlicesAreEqual(t testing.TB, got, want Notes) {
	t.Helper()
	sort.Slice(got, func(i, j int) bool { return got[i].Note < got[j].Note })
	sort.Slice(want, func(i, j int) bool { return want[i].Note < want[j].Note })
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`got = %v; want %v`, got, want)
	}
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

func newGetAllNotesRequest(t testing.TB) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/notes", nil)
	if err != nil {
		t.Fatalf("Unable to build request newGetAllNotesRequest %q", err)
	}
	return req
}

func getAllNotesFromResponse(t testing.TB, body io.Reader) (allNotes Notes) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&allNotes)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into map[int]Notes", err)
	}
	return
}

func getNotesByIdFromResponse(t testing.TB, body io.Reader) (notes Notes) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&notes)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into map[int]Notes", err)
	}
	return
}

func assertAllNotes(t testing.TB, got, wantedNotes Notes) {
	t.Helper()
	if !reflect.DeepEqual(got, wantedNotes) {
		t.Errorf("got %v want %v", got, wantedNotes)
	}
}

func assertNotesById(t testing.TB, got, want Notes) {
	t.Helper()
	assertSlicesHaveSameLength(t, got, want)
	assertStringSlicesAreEqual(t, got, want)
}

func assertAddNoteCalls(t testing.TB, got, want Notes) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func assertLoggerInfolnCalls(t testing.TB, got []any, want []string) {
	t.Helper()
	var gotStrings []string
	for _, val := range got {
		if s, ok := val.(string); ok {
			gotStrings = append(gotStrings, s)
		} else {
			t.Fatalf("Could not convert infolnCall to string")
		}
	}
	if !reflect.DeepEqual(gotStrings, want) {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}
