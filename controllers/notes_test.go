package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

type StubNotesStoreAddNoteErrors struct {
	StubNotesStore
}

func (sns *StubNotesStoreAddNoteErrors) AddNote(note Note) error {
	return errors.New("Error stub")
}

type fmtCallf struct {
	format string
	a      []any
}

type StubLogger struct {
	infofCalls []fmtCallf
	errorfCall []fmtCallf
}

func (sl *StubLogger) Infof(format string, a ...any) {
	sl.infofCalls = append(sl.infofCalls, fmtCallf{format: format, a: a})
}

func (sl *StubLogger) Errorf(format string, a ...any) {
	sl.errorfCall = append(sl.errorfCall, fmtCallf{format: format, a: a})
}

func (sl *StubLogger) Reset() {
	sl.infofCalls = []fmtCallf{}
	sl.errorfCall = []fmtCallf{}
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
		assertGotCallsEqualsWantCalls(t, logger.infofCalls, []fmtCallf{
			{format: "%s request to %s received", a: []any{"GET", "/notes"}},
		})
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
		assertGotCallsEqualsWantCalls(t, logger.infofCalls, []fmtCallf{
			{format: "%s request to %s received", a: []any{"GET", "/notes/1"}},
			{format: "%s request to %s received", a: []any{"GET", "/notes/2"}},
			{format: "%s request to %s received", a: []any{"GET", "/notes/100"}},
		})
	})

	t.Run("adds a note with POST", func(t *testing.T) {
		logger.Reset()
		note := NewNote(1, "Test note")
		request := newPostRequestWithNote(t, note, "/notes/1")
		response := httptest.NewRecorder()
		notesC.ProcessAddNote(response, request)

		wantAddNoteCalls := Notes{note}
		assertStatusCode(t, response.Result().StatusCode, http.StatusAccepted)
		assertAddNoteCalls(t, notesStore.addNoteCalls, wantAddNoteCalls)
		assertGotCallsEqualsWantCalls(t, logger.infofCalls, []fmtCallf{
			{format: "%s request to %s received", a: []any{"POST", "/notes/1"}},
		})
	})

	t.Run("test invalid json body", func(t *testing.T) {
		logger.Reset()
		badRequest := newPostRequestFromBody(t, "{}}", "/notes/1")
		response := httptest.NewRecorder()
		notesC.ProcessAddNote(response, badRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusBadRequest)
		assertRightErrorCall(t, logger.errorfCall[0], "%w: %w", ErrUnmarshalRequestBody)
	})

	t.Run("test invalid request body", func(t *testing.T) {
		logger.Reset()

		badRequest := newInvalidBodyPostRequest(t)
		response := httptest.NewRecorder()
		notesC.ProcessAddNote(response, badRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusBadRequest)
		assertRightErrorCall(t, logger.errorfCall[0], "%w: %w", ErrInvalidRequestBody)
	})

	t.Run("test AddNote returns error", func(t *testing.T) {
		logger.Reset()
		notesC := &NotesCtrlr{NotesStore: &StubNotesStoreAddNoteErrors{}, Logger: &logger}

		request := newPostRequestWithNote(t, NewNote(1, "Test note"), "/notes/1")
		response := httptest.NewRecorder()

		notesC.ProcessAddNote(response, request)
		assertStatusCode(t, response.Result().StatusCode, http.StatusInternalServerError)
		assertRightErrorCall(t, logger.errorfCall[0], "%w: %w", ErrDBResourceCreation)
	})

	t.Run("test false url parameters throws error", func(t *testing.T) {
		logger.Reset()
		response := httptest.NewRecorder()

		badUrlParam := "notAnInt"
		request, err := http.NewRequest(http.MethodGet, "", nil)
		assertNoError(t, err)

		badRequest := WithUrlParam(request, "id", fmt.Sprintf("%v", badUrlParam))
		notesC.GetNotesByID(response, badRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusBadRequest)
		assertRightErrorCall(t, logger.errorfCall[0], "%w: %w", ErrInvalidUserID)
	})
}

func newInvalidBodyPostRequest(t testing.TB) *http.Request {
	note := Note{UserID: 1, Note: "hello"}
	requestBody := map[string]Note{"wrong_key": note}
	buf := encodeRequestBodyAddNote(t, requestBody)
	badRequest, err := http.NewRequest(http.MethodPost, "", buf)
	assertNoError(t, err)
	return badRequest

}

func newPostRequestFromBody(t testing.TB, requestBody string, url string) *http.Request {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(requestBody)
	assertNoError(t, err)
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	assertNoError(t, err)
	return req
}

func encodeRequestBodyAddNote(t testing.TB, rb map[string]Note) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(rb)
	assertNoError(t, err)
	return buf
}

func newPostRequestWithNote(t testing.TB, note Note, url string) *http.Request {
	requestBody := map[string]Note{"note": note}
	buf := encodeRequestBodyAddNote(t, requestBody)
	request, err := http.NewRequest(http.MethodPost, url, buf)
	assertNoError(t, err)
	return request
}

func newGetNotesByUserIdRequest(t testing.TB, userID int) *http.Request {
	url := fmt.Sprintf("/notes/%v", userID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
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

func assertGotCallsEqualsWantCalls(t testing.TB, got, want []fmtCallf) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertRightErrorCall(t testing.TB, errorCall fmtCallf, wantFormat string, wantErr error) {
	t.Helper()
	gotFormat := errorCall.format
	if gotFormat != wantFormat {
		t.Errorf(`got = %v; want %v`, gotFormat, wantFormat)
	}
	if gotErr, ok := errorCall.a[0].(error); ok {
		if !errors.Is(gotErr, wantErr) {
			t.Errorf(`got = %v; want %v`, gotErr, wantErr)
		}
	} else {
		t.Errorf("Could not convert to error")
	}
}
