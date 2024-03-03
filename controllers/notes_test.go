package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"fmt"
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
		assertLoggingCalls(t, logger.infofCalls, []string{"Success: GetAllNotes"})
	})

	t.Run("Return notes for user with userID", func(t *testing.T) {
		logger.Reset()
		testCases := []struct {
			userID     int
			statusCode int
		}{
			{1, http.StatusOK},
			{2, http.StatusOK},
			{100, http.StatusNotFound},
		}

		for _, tc := range testCases {
			response := httptest.NewRecorder()
			request := newGetNotesByUserIdRequest(t, tc.userID)
			notesC.GetNotesByUserID(response, request)
			assertStatusCode(t, response.Result().StatusCode, tc.statusCode)
		}
		assertEqualIntSlice(t, notesStore.getNotesByUserIDCalls, []int{1, 2, 100})
		assertLoggingCalls(t, logger.infofCalls, []string{
			"Success: GetNotesByUserID with userID %d",
			"Success: GetNotesByUserID with userID %d",
		})
		assertLoggingCalls(t, logger.errorfCall, []string{"GetNotesByUserID user not Found: %w"})
	})

	t.Run("test false url parameters throw error", func(t *testing.T) {
		logger.Reset()

		badID := "notAnInt"
		badRequest := newRequestWithBadIdParam(t, badID)
		response := httptest.NewRecorder()
		notesC.GetNotesByUserID(response, badRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusBadRequest)
		assertLoggingCalls(t, logger.errorfCall, []string{"GetNotesByUserID invalid userID: %w"})
	})

	t.Run("POST a Note", func(t *testing.T) {
		logger.Reset()
		userID, note := 1, "Test note"

		request := newPostRequestWithNote(t, note)
		request = WithUrlParam(request, "userID", fmt.Sprintf("%d", userID))
		response := httptest.NewRecorder()
		notesC.ProcessAddNote(response, request)

		wantAddNoteCalls := []AddNoteCall{{userID: userID, note: note}}
		assertStatusCode(t, response.Result().StatusCode, http.StatusAccepted)
		assertSlicesAnyAreEqual(t, notesStore.addNoteCalls, wantAddNoteCalls)
		assertLoggingCalls(t, logger.infofCalls, []string{"Success: ProcessAddNote with userID %d and note %v"})
	})

	t.Run("test invalid json body", func(t *testing.T) {
		logger.Reset()
		badRequest := newPostRequestFromBody(t, "{}}")
		badRequest = WithUrlParam(badRequest, "userID", fmt.Sprintf("%d", 1))
		response := httptest.NewRecorder()
		notesC.ProcessAddNote(response, badRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusBadRequest)
		assertLoggingCalls(t, logger.errorfCall, []string{"ProcessAddNote invalid json: %w"})
	})

	t.Run("test invalid request body", func(t *testing.T) {
		logger.Reset()

		badRequest := newInvalidBodyPostRequest(t)
		badRequest = WithUrlParam(badRequest, "userID", fmt.Sprintf("%d", 1))
		response := httptest.NewRecorder()
		notesC.ProcessAddNote(response, badRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusBadRequest)
		assertLoggingCalls(t, logger.errorfCall, []string{"ProcessAddNote invalid body: %w"})
	})

	t.Run("test AddNote and Note already present", func(t *testing.T) {
		logger.Reset()

		request := newPostRequestWithNote(t, "Note already present")
		request = WithUrlParam(request, "userID", fmt.Sprintf("%d", 1))
		response := httptest.NewRecorder()

		notesC.ProcessAddNote(response, request)
		assertStatusCode(t, response.Result().StatusCode, http.StatusInternalServerError)
		assertLoggingCalls(t, logger.errorfCall, []string{"ProcessAddNote DBerror: %w"})
	})

	t.Run("Delete a Note", func(t *testing.T) {
		logger.Reset()

		userID, noteID := 1, 2
		request, err := http.NewRequest(http.MethodDelete, "", nil)
		assertNoError(t, err)
		request = WithUrlParams(request, map[string]string{
			"userID": fmt.Sprintf("%d", userID),
			"noteID": fmt.Sprintf("%d", noteID),
		})

		response := httptest.NewRecorder()
		notesC.Delete(response, request)
		wantDeleteNoteCalls := []DeleteCall{{userID: userID, noteID: noteID}}
		assertStatusCode(t, response.Result().StatusCode, http.StatusNoContent)
		assertSlicesAnyAreEqual(t, notesStore.deleteNoteCalls, wantDeleteNoteCalls)
		assertLoggingCalls(t, logger.errorfCall, []string{"Success: Delete noteID %v userID %v"})
	})

	t.Run("Deletion fail", func(t *testing.T) {
		logger.Reset()

		userID, noteID := 50, 50
		request, err := http.NewRequest(http.MethodDelete, "", nil)
		assertNoError(t, err)
		request = WithUrlParams(request, map[string]string{
			"userID": fmt.Sprintf("%d", userID),
			"noteID": fmt.Sprintf("%d", noteID),
		})

		response := httptest.NewRecorder()
		notesC.Delete(response, request)

		assertStatusCode(t, response.Result().StatusCode, http.StatusNotFound)
		assertLoggingCalls(t, logger.errorfCall, []string{"Delete DBError: %w"})
	})

	t.Run("Edit a Note", func(t *testing.T) {
		logger.Reset()

		userID, noteID, note := 1, 1, "New note text"
		putRequest := newPutRequestWithNote(t, "New note text")
		WithUrlParams(putRequest, map[string]string{
			"userID": fmt.Sprintf("%d", userID),
			"noteID": fmt.Sprintf("%d", noteID),
		})
		response := httptest.NewRecorder()
		notesC.Edit(response, putRequest)

		wantEditCalls := []EditCall{{userID: userID, noteID: noteID, note: note}}
		assertStatusCode(t, response.Result().StatusCode, http.StatusOK)
		assertSlicesAnyAreEqual(t, notesStore.editNoteCalls, wantEditCalls)

		// wantEditNoteCalls := Notes{note}
		// assertAddNoteCallsEqual(t, notesStore.editNoteCalls, wantEditNoteCalls)
		// assertLoggingCalls(t, logger.infofCalls, []fmtCallf{
		// 	{format: "%s request to %s received", a: []any{"PUT", "/notes/1"}},
		// })
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

func WithUrlParams(r *http.Request, params map[string]string) *http.Request {
	chiCtx := chi.NewRouteContext()
	req := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	for key, value := range params {
		chiCtx.URLParams.Add(key, value)
	}
	return req
}

func encodeRequestBodyAddNote(t testing.TB, rb map[string]string) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(rb)
	assertNoError(t, err)
	return buf
}

func newGetNotesByUserIdRequest(t testing.TB, userID int) *http.Request {
	url := fmt.Sprintf("/users/%v/notes", userID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Could not build request newPostAddNoteRequest: %q", err)
	}
	return WithUrlParam(request, "userID", fmt.Sprintf("%v", userID))
}

func newGetAllNotesRequest(t testing.TB) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/notes", nil)
	if err != nil {
		t.Fatalf("Unable to build request newGetAllNotesRequest %q", err)
	}
	return req
}

func newDeleteRequest(t testing.TB, userID, noteID int) *http.Request {
	request, err := http.NewRequest(http.MethodDelete, "", nil)
	assertNoError(t, err)
	request = WithUrlParam(request, "id", fmt.Sprintf("%d", noteID))
	return request
}

func newPostRequestWithNote(t testing.TB, note string) *http.Request {
	requestBody := map[string]string{"note": note}
	buf := encodeRequestBodyAddNote(t, requestBody)
	request, err := http.NewRequest(http.MethodPost, "", buf)
	assertNoError(t, err)
	return request
}

func newPutRequestWithNote(t testing.TB, note string) *http.Request {
	requestBody := map[string]string{"note": note}
	buf := encodeRequestBodyAddNote(t, requestBody)
	request, err := http.NewRequest(http.MethodPut, "", buf)
	assertNoError(t, err)
	return request
}

func newPostRequestFromBody(t testing.TB, requestBody string) *http.Request {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(requestBody)
	assertNoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "", &buf)
	assertNoError(t, err)
	return req
}

func newRequestWithBadIdParam(t testing.TB, badID string) *http.Request {
	request, err := http.NewRequest(http.MethodGet, "", nil)
	assertNoError(t, err)
	return WithUrlParam(request, "id", fmt.Sprintf("%v", badID))
}

func newInvalidBodyPostRequest(t testing.TB) *http.Request {
	requestBody := map[string]string{"wrong_key": "some text"}
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

func assertSlicesAnyAreEqual[T any](t testing.TB, gotSlice, wantSlice []T) {
	t.Helper()
	assertLengthSlice(t, gotSlice, len(wantSlice))
	for _, want := range wantSlice {
		found := false
		for _, got := range gotSlice {
			if reflect.DeepEqual(got, want) {
				found = true
			}
		}
		if !found {
			t.Errorf("want %v not found in gotSlice %v", want, gotSlice)
		}
	}
}

func assertDeleteNoteCallsEqual(t testing.TB, gotCalls, wantCalls []DeleteCall) {
	t.Helper()
	assertLengthSlice(t, gotCalls, len(wantCalls))
	for _, want := range wantCalls {
		found := false
		for _, got := range gotCalls {
			if reflect.DeepEqual(got, want) {
				found = true
			}
		}
		if !found {
			t.Errorf("want %v not found in gotCalls %v", want, gotCalls)
		}
	}
}

func assertLoggingCalls(t testing.TB, got, want []string) {
	t.Helper()
	assertSlicesAnyAreEqual(t, got, want)
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
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
