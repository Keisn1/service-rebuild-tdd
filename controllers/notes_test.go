package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"fmt"

	"errors"
	ctrls "github.com/Keisn1/note-taking-app/controllers"
	"github.com/Keisn1/note-taking-app/domain"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetAllNotes(t *testing.T) {
	notes := domain.Notes{
		{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
		{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
		{NoteID: 3, UserID: 2, Note: "Note 1 user 2"},
		{NoteID: 4, UserID: 2, Note: "Note 2 user 2"},
	}

	mNotesStore := &mockNotesStore{}
	mLogger := &mockLogger{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore, mLogger)

	testCases := []struct {
		name         string
		handler      http.HandlerFunc
		wantBody     string
		setupMock    func(*testing.T)
		setupRequest func(*testing.T) *http.Request
		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantBody string)
	}{
		{
			name:    "GetAllNotes Happy path",
			handler: notesCtrl.GetAllNotes,
			wantBody: mustEncode(t, domain.Notes{
				{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
				{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
				{NoteID: 3, UserID: 2, Note: "Note 1 user 2"},
				{NoteID: 4, UserID: 2, Note: "Note 2 user 2"},
			}),
			setupMock: func(t *testing.T) {
				t.Helper()
				resetMocks(mNotesStore, mLogger)
				mNotesStore.On("GetAllNotes").Return(notes, nil)
				mLogger.On("Infof", "Success: GetAllNotes").Return(nil)
			},
			setupRequest: func(t *testing.T) *http.Request {
				return httptest.NewRequest("GET", "/users/notes", nil)
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantBody string) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, rr.Body.String(), wantBody)
				mNotesStore.AssertCalled(t, "GetAllNotes")
				mLogger.AssertCalled(t, "Infof", "Success: GetAllNotes")
			},
		},
		{
			name:    "GetAllNotes Error DB",
			handler: notesCtrl.GetAllNotes,
			setupMock: func(t *testing.T) {
				resetMocks(mNotesStore, mLogger)
				mNotesStore.On("GetAllNotes").Return(domain.Notes{}, errors.New("error notesStore.GetAllNotes"))
				mLogger.On("Errorf", "%s: %w", "GetAllNotes", errors.New("error notesStore.GetAllNotes")).Return(nil)
			},
			setupRequest: func(t *testing.T) *http.Request {
				t.Helper()
				return httptest.NewRequest("GET", "/users/notes", nil)
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantBody string) {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
				mNotesStore.AssertCalled(t, "GetAllNotes")
				mLogger.AssertCalled(t, "Errorf", "%s: %w", "GetAllNotes", errors.New("error notesStore.GetAllNotes"))
			},
		},
	}

	for _, tc := range testCases {
		tc.setupMock(t)
		req := tc.setupRequest(t)
		rr := httptest.NewRecorder()
		tc.handler(rr, req)
		tc.assertions(t, rr, tc.wantBody)
	}
}

func TestGetNoteByUserIDandNoteID(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	mLogger := &mockLogger{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore, mLogger)

	testCases := []struct {
		name         string
		handler      http.HandlerFunc
		wantBody     string
		setupMock    func(*testing.T)
		setupRequest func(*testing.T) *http.Request
		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantBody string)
	}{
		{
			name:     "GetNoteByUserIDandNoteID success",
			handler:  notesCtrl.GetNoteByUserIDAndNoteID,
			wantBody: mustEncode(t, domain.Notes{{NoteID: 1, UserID: 1, Note: "Note 1 user 1"}}),
			setupMock: func(t *testing.T) {
				resetMocks(mNotesStore, mLogger)
				userID, noteID := 1, 1
				note := fmt.Sprintf("Note %v user %v", userID, noteID)
				mNotesStore.On("GetNoteByUserIDAndNoteID", userID, noteID).Return(domain.Notes{{NoteID: noteID, UserID: userID, Note: note}}, nil)
				mLogger.On("Infof", "Success: GetNoteByUserIDAndNoteID with userID %v and noteID %v", userID, noteID).Return(nil)
			},
			setupRequest: func(t *testing.T) *http.Request {
				t.Helper()
				userID, noteID := 1, 1
				req := httptest.NewRequest("GET", "/users/{userID}/notes/{noteID}", nil)
				return WithUrlParams(req, Params{
					"userID": strconv.Itoa(userID),
					"noteID": strconv.Itoa(noteID),
				})
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantBody string) {
				userID, noteID := 1, 1
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertCalled(t, "GetNoteByUserIDAndNoteID", userID, noteID)
				mLogger.AssertCalled(t, "Infof", "Success: GetNoteByUserIDAndNoteID with userID %v and noteID %v", userID, noteID)
			},
		},
		{
			name:    "GetNoteByUserIDandNoteID invalid userID",
			handler: notesCtrl.GetNoteByUserIDAndNoteID,
			setupMock: func(t *testing.T) {
				resetMocks(mNotesStore, mLogger)
				userID := -1
				mLogger.On("Errorf", "GetNoteByUserIDandNoteID: invalid userID %v", userID).Return(nil)
			},
			setupRequest: func(t *testing.T) *http.Request {
				t.Helper()
				userID, noteID := -1, 1
				req := httptest.NewRequest("GET", "/users/{userID}/notes/{noteID}", nil)
				return WithUrlParams(req, Params{
					"userID": strconv.Itoa(userID),
					"noteID": strconv.Itoa(noteID),
				})
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantBody string) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
				mNotesStore.AssertNotCalled(t, "GetNoteByUserIDandNoteID")
				mLogger.AssertCalled(t, "Errorf", "GetNoteByUserIDandNoteID: invalid userID %v", -1)
			},
		},
		{
			name:    "GetNoteByUserIDandNoteID invalid userID",
			handler: notesCtrl.GetNoteByUserIDAndNoteID,
			setupMock: func(t *testing.T) {
				resetMocks(mNotesStore, mLogger)
				noteID := -1
				mLogger.On("Errorf", "GetNoteByUserIDandNoteID: invalid noteID %v", noteID).Return(nil)
			},
			setupRequest: func(t *testing.T) *http.Request {
				t.Helper()
				userID, noteID := 1, -1
				req := httptest.NewRequest("GET", "/users/{userID}/notes/{noteID}", nil)
				return WithUrlParams(req, Params{
					"userID": strconv.Itoa(userID),
					"noteID": strconv.Itoa(noteID),
				})
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantBody string) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
				mNotesStore.AssertNotCalled(t, "GetNoteByUserIDandNoteID")
				mLogger.AssertCalled(t, "Errorf", "GetNoteByUserIDandNoteID: invalid noteID %v", -1)
			},
		},
		{
			name:    "GetNoteByUserIDandNoteID DBError",
			handler: notesCtrl.GetNoteByUserIDAndNoteID,
			setupMock: func(t *testing.T) {
				resetMocks(mNotesStore, mLogger)
				userID, noteID := 1, 1
				mNotesStore.On("GetNoteByUserIDAndNoteID", userID, noteID).Return(domain.Notes{}, errors.New("error notesStore.GetNoteByUserIDAndNoteID"))
				mLogger.On("Errorf", "%s: %w", "GetNoteByUserIDAndNoteID userID 1 and noteID 1", errors.New("error notesStore.GetNoteByUserIDAndNoteID")).Return(nil)
			},
			setupRequest: func(t *testing.T) *http.Request {
				t.Helper()
				userID, noteID := 1, 1
				req := httptest.NewRequest("GET", "/users/{userID}/notes/{noteID}", nil)
				return WithUrlParams(req, Params{
					"userID": strconv.Itoa(userID),
					"noteID": strconv.Itoa(noteID),
				})
			},

			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantBody string) {
				userID, noteID := 1, 1
				assert.Equal(t, http.StatusNotFound, rr.Code)
				mNotesStore.AssertCalled(t, "GetNoteByUserIDAndNoteID", userID, noteID)
				mLogger.AssertCalled(t, "Errorf", "%s: %w", "GetNoteByUserIDAndNoteID userID 1 and noteID 1", errors.New("error notesStore.GetNoteByUserIDAndNoteID"))
			},
		},
	}

	for _, tc := range testCases {
		tc.setupMock(t)
		req := tc.setupRequest(t)
		rr := httptest.NewRecorder()
		tc.handler(rr, req)
		tc.assertions(t, rr, tc.wantBody)
	}
}

func TestNotes2(t *testing.T) {
	notes := domain.Notes{
		{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
		{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
		{NoteID: 3, UserID: 2, Note: "Note 1 user 2"},
		{NoteID: 4, UserID: 2, Note: "Note 2 user 2"},
	}

	notesStore := NewStubNotesStore(notes)
	logger := NewStubLogger()
	notesCtrl := ctrls.NewNotesCtrlr(notesStore, logger)

	t.Run("Return notes for user with userID", func(t *testing.T) {
		logger.Reset()
		testCases := []struct {
			userID     int
			wantNotes  domain.Notes
			statusCode int
		}{
			{1, domain.Notes{
				{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
				{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
			}, http.StatusOK},
			{2, domain.Notes{
				{NoteID: 3, UserID: 2, Note: "Note 1 user 2"},
				{NoteID: 4, UserID: 2, Note: "Note 2 user 2"},
			}, http.StatusOK}, // no notes for user with userID 2
			{-1, domain.Notes{}, http.StatusInternalServerError}, // simulating DBError
		}

		for _, tc := range testCases {
			response := httptest.NewRecorder()
			request := newGetNotesByUserIdRequest(t, tc.userID)
			notesCtrl.GetNotesByUserID(response, request)
			if tc.userID != -1 {
				gotNotes := decodeBodyNotes(t, response.Body)
				assertSlicesAnyAreEqual(t, gotNotes, tc.wantNotes)
			}

			assertStatusCode(t, response.Result().StatusCode, tc.statusCode)
		}
		assertEqualIntSlice(t, notesStore.getNotesByUserIDCalls, []int{1, 2, -1})
		assertLoggingCalls(t, logger.infofCalls, []string{
			"Success: GetNotesByUserID with userID 1",
			"Success: GetNotesByUserID with userID 2",
		})
		assertLoggingCalls(t, logger.errorfCall, []string{
			fmt.Sprintf("GetNotesByUserID userID -1 %v", ctrls.ErrDB.Error()),
		})
	})

	t.Run("test false url parameters throw error", func(t *testing.T) {
		logger.Reset()

		badID := "notAnInt"
		badRequest := newRequestWithBadIdParam(t, badID)
		response := httptest.NewRecorder()
		notesCtrl.GetNotesByUserID(response, badRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusBadRequest)
		assertLoggingCalls(t, logger.errorfCall, []string{"GetNotesByUserID invalid userID:"})
	})

	t.Run("POST a Note", func(t *testing.T) {
		logger.Reset()
		userID, note := 1, "Test note"

		request := newPostRequestWithNoteAndUrlParam(t, note, "userID", fmt.Sprintf("%d", userID))
		response := httptest.NewRecorder()
		notesCtrl.Add(response, request)

		wantAddNoteCalls := []AddNoteCall{{userID: userID, note: note}}
		assertStatusCode(t, response.Result().StatusCode, http.StatusAccepted)
		assertSlicesAnyAreEqual(t, notesStore.addNoteCalls, wantAddNoteCalls)
		assertLoggingCalls(t, logger.infofCalls, []string{"Success: ProcessAddNote with userID 1 and note Test note"})
	})

	t.Run("test invalid json body", func(t *testing.T) {
		logger.Reset()
		badRequest := newPostRequestFromBody(t, "{}}")
		badRequest = WithUrlParam(badRequest, "userID", "1")
		response := httptest.NewRecorder()
		notesCtrl.Add(response, badRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusBadRequest)
		assertLoggingCalls(t, logger.errorfCall, []string{"ProcessAddNote invalid json:"})
	})

	t.Run("test invalid request body", func(t *testing.T) {
		logger.Reset()

		badRequest := newInvalidBodyPostRequest(t)
		badRequest = WithUrlParam(badRequest, "userID", fmt.Sprintf("%d", 1))
		response := httptest.NewRecorder()
		notesCtrl.Add(response, badRequest)

		assertStatusCode(t, response.Result().StatusCode, http.StatusBadRequest)
		assertLoggingCalls(t, logger.errorfCall, []string{"ProcessAddNote invalid body:"})
	})

	// t.Run("test AddNote and Note already present", func(t *testing.T) {
	// 	logger.Reset()

	// 	request := newPostRequestWithNote(t, "Note already present")
	// 	request = WithUrlParam(request, "userID", fmt.Sprintf("%d", 1))
	// 	response := httptest.NewRecorder()

	// 	notesCtrl.Add(response, request)
	// 	assertStatusCode(t, response.Result().StatusCode, http.StatusConflict)
	// 	assertLoggingCalls(t, logger.errorfCall, []string{"ProcessAddNote DBerror:"})
	// })

	t.Run("Delete a Note", func(t *testing.T) {
		logger.Reset()

		userID, noteID := 1, 2
		request, err := http.NewRequest(http.MethodDelete, "", nil)
		assertNoError(t, err)
		request = WithUrlParams(request, Params{
			"userID": strconv.Itoa(userID),
			"noteID": strconv.Itoa(noteID),
		})

		response := httptest.NewRecorder()
		notesCtrl.Delete(response, request)
		wantDeleteNoteCalls := []DeleteCall{{userID: userID, noteID: noteID}}
		assertStatusCode(t, response.Result().StatusCode, http.StatusNoContent)
		assertSlicesAnyAreEqual(t, notesStore.deleteNoteCalls, wantDeleteNoteCalls)
		assertLoggingCalls(t, logger.infofCalls, []string{"Success: Delete noteID 2 userID 1"})
	})

	// t.Run("Deletion fail", func(t *testing.T) {
	// 	logger.Reset()

	// 	userID, noteID := 50, 50
	// 	request, err := http.NewRequest(http.MethodDelete, "", nil)
	// 	assertNoError(t, err)
	// 	request = WithUrlParams(request, Params{
	// 		"userID": strconv.Itoa(userID),
	// 		"noteID": strconv.Itoa(noteID),
	// 	})

	// 	response := httptest.NewRecorder()
	// 	notesCtrl.Delete(response, request)

	// 	assertStatusCode(t, response.Result().StatusCode, http.StatusInternalServerError)
	// 	assertLoggingCalls(t, logger.errorfCall, []string{"Delete DBerror:"})
	// })

	t.Run("Edit a Note", func(t *testing.T) {
		logger.Reset()

		userID, noteID, note := 1, 1, "New note text"
		putRequest := newPutRequestWithNoteAndUrlParams(t, "New note text", Params{
			"userID": strconv.Itoa(userID),
			"noteID": strconv.Itoa(noteID),
		})
		response := httptest.NewRecorder()
		notesCtrl.Edit(response, putRequest)

		wantEditCalls := []EditCall{{userID: userID, noteID: noteID, note: note}}
		assertStatusCode(t, response.Result().StatusCode, http.StatusOK)
		assertSlicesAnyAreEqual(t, notesStore.editNoteCalls, wantEditCalls)
		assertLoggingCalls(t, logger.infofCalls, []string{
			"Success: Edit: userID 1 noteID 1 note New note text",
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

type Params map[string]string

func WithUrlParams(r *http.Request, params Params) *http.Request {
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
func newGetNoteByUserIDAndNoteIDRequest(t testing.TB, userID, noteID int) *http.Request {
	request, err := http.NewRequest(http.MethodGet, "", nil)
	assertNoError(t, err)
	return WithUrlParams(request, Params{
		"userID": strconv.Itoa(userID),
		"noteID": strconv.Itoa(noteID),
	})
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
	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("Unable to build request newGetAllNotesRequest %q", err)
	}
	return req
}

func newPostRequestWithNoteAndUrlParam(t testing.TB, note string, key, value string) *http.Request {
	request := newPostRequestWithNote(t, note)
	request = WithUrlParam(request, key, value)
	return request
}

func newPostRequestWithNote(t testing.TB, note string) *http.Request {
	requestBody := map[string]string{"note": note}
	buf := encodeRequestBodyAddNote(t, requestBody)
	request, err := http.NewRequest(http.MethodPost, "", buf)
	assertNoError(t, err)
	return request
}

func newPutRequestWithNoteAndUrlParams(t testing.TB, note string, params Params) *http.Request {
	request := newPutRequestWithNote(t, note)
	request = WithUrlParams(request, params)
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

func assertSlicesSameLength[T any](t testing.TB, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf(`len(got) = %v; len(want) %v`, len(got), len(want))
	}
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func assertSlicesAnyAreEqual[T any](t testing.TB, gotSlice, wantSlice []T) {
	t.Helper()
	assertSlicesSameLength(t, gotSlice, wantSlice)
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

func assertLoggingCalls(t testing.TB, gotCalls, wantCalls []string) {
	t.Helper()
	for _, want := range wantCalls {
		found := false
		for _, got := range gotCalls {
			if strings.HasPrefix(got, want) {
				found = true
			}
		}
		if !found {
			t.Errorf("want %v - wasn't found in gotCalls %v", want, gotCalls)
		}
	}
}

func assertGetAllNotesGotCalled(t testing.TB, allNotesGotCalled bool) {
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

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func resetMocks(mNotesStore *mockNotesStore, mLogger *mockLogger) {
	mNotesStore.Reset()
	mLogger.Reset()
}

func mustEncode(t *testing.T, a any) string {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(a); err != nil {
		t.Fatalf("encoding json: %v", err)
	}
	return buf.String()
}
