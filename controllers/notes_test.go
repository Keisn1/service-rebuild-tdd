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
	"testing"

	"github.com/go-chi/chi"
)

func TestNotes(t *testing.T) {
	notesStore := NewStubNotesStore()
	logger := NewStubLogger()
	notesC := NewNotesCtrlr(notesStore, logger)

	t.Run("Server returns all Notes", func(t *testing.T) {
		logger.Reset()

		request := newGetAllNotesRequest(t)
		response := httptest.NewRecorder()
		notesC.GetAllNotes(response, request)

		assertStatusCode(t, response.Result().StatusCode, http.StatusOK)
		assertAllNotesGotCalled(t, notesStore.allNotesGotCalled)
		assertLoggingCalls(t, logger.infofCalls, []fmtCallf{
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
			notesC.GetNotesByUserID(response, request)

			assertStatusCode(t, response.Result().StatusCode, tc.statusCode)
			got := getNotesFromResponse(t, response.Body)
			assertNotesEqual(t, got, tc.want)
		}
		assertEqualIntSlice(t, notesStore.getNotesByUserIDCalls, []int{1, 2, 3})
		assertLoggingCalls(t, logger.infofCalls, []fmtCallf{
			{format: "%s request to %s received", a: []any{"GET", "/notes/1"}},
			{format: "%s request to %s received", a: []any{"GET", "/notes/2"}},
			{format: "%s request to %s received", a: []any{"GET", "/notes/100"}},
		})
	})

	t.Run("POST a Note", func(t *testing.T) {
		logger.Reset()
		note := NewNote(1, "Test note")
		request := newPostRequestWithNote(t, note, "/notes/1")
		response := httptest.NewRecorder()
		notesC.ProcessAddNote(response, request)

		wantAddNoteCalls := Notes{note}
		assertStatusCode(t, response.Result().StatusCode, http.StatusAccepted)
		assertNotesEqual(t, notesStore.addNoteCalls, wantAddNoteCalls)
		assertLoggingCalls(t, logger.infofCalls, []fmtCallf{
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

		// note already present
		request := newPostRequestWithNote(t, NewNote(1, "Note 1 user 1"), "/notes/1")
		response := httptest.NewRecorder()

		notesC.ProcessAddNote(response, request)
		assertStatusCode(t, response.Result().StatusCode, http.StatusInternalServerError)
		assertRightErrorCall(t, logger.errorfCall[0], "%w: %w", ErrDBResourceCreation)
	})

	t.Run("test false url parameters throws error", func(t *testing.T) {
		logger.Reset()

		badRequest := newRequestWithBadIdParam(t)
		response := httptest.NewRecorder()
		notesC.GetNotesByUserID(response, badRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusBadRequest)
		assertRightErrorCall(t, logger.errorfCall[0], "%w: %w", ErrInvalidUserID)
	})

	t.Run("Delete a Note", func(t *testing.T) {
		logger.Reset()

		deleteRequest := newDeleteRequest(t, 1)
		response := httptest.NewRecorder()
		notesC.Delete(response, deleteRequest)

		wantDeleteNoteCalls := []int{1}
		assertStatusCode(t, response.Result().StatusCode, http.StatusNoContent)
		assertEqualIntSlice(t, notesStore.deleteNoteCalls, wantDeleteNoteCalls)
	})

	t.Run("Deletion fail", func(t *testing.T) {
		logger.Reset()

		deleteRequest := newDeleteRequest(t, 50) // id not present
		response := httptest.NewRecorder()
		notesC.Delete(response, deleteRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusNotFound)
		assertRightErrorCall(t, logger.errorfCall[0], "%w: %w", ErrDBResourceDeletion)
	})

	t.Run("Edit a Note", func(t *testing.T) {
		logger.Reset()
		notesStore := NewStubNotesStore()
		logger := NewStubLogger()
		notesC := NewNotesCtrlr(notesStore, logger)

		note := NewNote(1, "Edited note")
		putRequest := newPutRequestWithNote(t, note, "/notes/1")
		response := httptest.NewRecorder()
		notesC.Edit(response, putRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusOK)

		wantEditNoteCalls := Notes{note}
		assertNotesEqual(t, notesStore.editNoteCalls, wantEditNoteCalls)
		assertLoggingCalls(t, logger.infofCalls, []fmtCallf{
			{format: "%s request to %s received", a: []any{"PUT", "/notes/1"}},
		})
	})
}

// WithUrlParam returns a pointer to a request object with the given URL params
// added to a new chi.Context object.
func WithUrlParam(r *http.Request, key, value string) *http.Request {
	chiCtx := chi.NewRouteContext()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add(key, value)
	return r
}

func getNotesFromResponse(t testing.TB, body io.Reader) (notes Notes) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&notes)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into map[int]Notes", err)
	}
	return
}

func encodeRequestBodyAddNote(t testing.TB, rb map[string]Note) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(rb)
	assertNoError(t, err)
	return buf
}

func newGetNotesByUserIdRequest(t testing.TB, userID int) *http.Request {
	url := fmt.Sprintf("/notes/%v", userID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Could not build request newPostAddNoteRequest: %q", err)
	}
	return WithUrlParam(request, "id", fmt.Sprintf("%v", userID))
}

func newGetAllNotesRequest(t testing.TB) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/notes", nil)
	if err != nil {
		t.Fatalf("Unable to build request newGetAllNotesRequest %q", err)
	}
	return req
}

func newDeleteRequest(t testing.TB, id int) *http.Request {
	url := fmt.Sprintf("/notes/%v", id)
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	assertNoError(t, err)
	request = WithUrlParam(request, "id", fmt.Sprintf("%d", id))
	return request
}

func newPostRequestWithNote(t testing.TB, note Note, url string) *http.Request {
	requestBody := map[string]Note{"note": note}
	buf := encodeRequestBodyAddNote(t, requestBody)
	request, err := http.NewRequest(http.MethodPost, url, buf)
	assertNoError(t, err)
	return request
}

func newPostRequestFromBody(t testing.TB, requestBody string, url string) *http.Request {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(requestBody)
	assertNoError(t, err)
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	assertNoError(t, err)
	return req
}

func newPutRequestWithNote(t testing.TB, note Note, url string) *http.Request {
	requestBody := map[string]Note{"note": note}
	buf := encodeRequestBodyAddNote(t, requestBody)
	request, err := http.NewRequest(http.MethodPut, url, buf)
	assertNoError(t, err)
	return request
}

func newRequestWithBadIdParam(t testing.TB) *http.Request {
	badUrlParam := "notAnInt"
	request, err := http.NewRequest(http.MethodGet, "", nil)
	assertNoError(t, err)
	return WithUrlParam(request, "id", fmt.Sprintf("%v", badUrlParam))
}

func newInvalidBodyPostRequest(t testing.TB) *http.Request {
	note := Note{UserID: 1, Note: "hello"}
	requestBody := map[string]Note{"wrong_key": note}
	buf := encodeRequestBodyAddNote(t, requestBody)
	badRequest, err := http.NewRequest(http.MethodPost, "", buf)
	assertNoError(t, err)
	return badRequest

}

func assertLengthSlice[T any](t testing.TB, elements []T, want int) {
	t.Helper()
	if len(elements) != want {
		t.Errorf(`got = %v; want %v`, len(elements), want)
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

func assertNotesEqual(t testing.TB, gotNotes, wantedNotes Notes) {
	t.Helper()
	assertLengthSlice(t, gotNotes, len(wantedNotes))
	for _, want := range wantedNotes {
		found := false
		for _, got := range gotNotes {
			if reflect.DeepEqual(got, want) {
				found = true
			}
		}
		if !found {
			t.Errorf("want %v not found in gotNotes %v", want, gotNotes)
		}
	}
}

func assertLoggingCalls(t testing.TB, got, want []fmtCallf) {
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

func assertAllNotesGotCalled(t testing.TB, allNotesGotCalled bool) {
	t.Helper()
	if !allNotesGotCalled {
		t.Error("notesStore.AllNotes did not get called")
	}
}
func assertEqualIntSlice(t testing.TB, got, want []int) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}
