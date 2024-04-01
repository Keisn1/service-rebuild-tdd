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
	"io"
)

type urlParams struct {
	userID string
	noteID string
}

func TestGetAllNotes(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	mLogger := &mockLogger{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore, mLogger)

	notes := domain.Notes{
		{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
		{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
		{NoteID: 3, UserID: 2, Note: "Note 1 user 2"},
		{NoteID: 4, UserID: 2, Note: "Note 2 user 2"},
	}
	testCases := []struct {
		name         string
		handler      http.HandlerFunc
		mockNSParams mockNotesStoreParams
		wantStatus   int
		wantBody     string
		wantLogging  mockLoggingParams
		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams)
	}{
		{
			name:         "GetAllNotes Success",
			handler:      notesCtrl.GetAllNotes,
			mockNSParams: mockNotesStoreParams{method: "GetAllNotes", returnArguments: []any{notes, nil}},
			wantStatus:   http.StatusOK,
			wantBody:     mustEncode(t, notes).String(),
			wantLogging: mockLoggingParams{
				method:    "Infof",
				arguments: []any{"Success: GetAllNotes"},
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, rr.Body.String(), wantBody)
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				mLogger.AssertCalled(t, wL.method, wL.arguments...)
			},
		},
		{
			name:    "GetAllNotes Error DB",
			handler: notesCtrl.GetAllNotes,
			mockNSParams: mockNotesStoreParams{
				method:          "GetAllNotes",
				returnArguments: []any{domain.Notes{}, errors.New("error notesStore.GetAllNotes")},
			},
			wantBody:   "\n",
			wantStatus: http.StatusInternalServerError,
			wantLogging: mockLoggingParams{
				method:    "Errorf",
				arguments: []any{"%s: %w", "GetAllNotes", errors.New("error notesStore.GetAllNotes")},
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, rr.Body.String(), wantBody)
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				mLogger.AssertCalled(t, wL.method, wL.arguments...)
			},
		},
	}

	for _, tc := range testCases {
		mNotesStore.Setup(tc.mockNSParams)
		mLogger.Setup(tc.wantLogging)
		req := setupRequest(t, "/users/notes", urlParams{}, &bytes.Reader{})
		rr := httptest.NewRecorder()
		notesCtrl.GetAllNotes(rr, req)
		tc.assertions(t, rr, tc.wantStatus, tc.wantBody, tc.wantLogging, tc.mockNSParams)
	}
}

func TestGetNoteByUserIDandNoteID(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	mLogger := &mockLogger{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore, mLogger)

	testCases := []struct {
		name         string
		handler      http.HandlerFunc
		urlParams    urlParams
		mockNSParams func(urlP urlParams) mockNotesStoreParams
		wantStatus   int
		wantBody     func(urlParams) string
		wantLogging  func(urlParams) mockLoggingParams
		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams)
	}{
		{
			name:      "GetNoteByUserIDandNoteID success",
			handler:   notesCtrl.GetNoteByUserIDAndNoteID,
			urlParams: urlParams{userID: "1", noteID: "1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams {
				userID, noteID := mustConvUrlParamsToInt(t, up)
				note := fmt.Sprintf("Note %d user %d", userID, noteID)
				return mockNotesStoreParams{
					method:    "GetNoteByUserIDAndNoteID",
					arguments: []any{userID, noteID},
					returnArguments: []any{
						domain.Notes{{UserID: userID, NoteID: noteID, Note: note}},
						nil,
					},
				}
			},
			wantStatus: http.StatusOK,
			wantBody: func(up urlParams) string {
				userID, noteID := mustConvUrlParamsToInt(t, up)
				note := fmt.Sprintf("Note %d user %d", userID, noteID)
				return mustEncode(t, domain.Notes{{UserID: userID, NoteID: noteID, Note: note}}).String()
			},
			wantLogging: func(up urlParams) mockLoggingParams {
				return mockLoggingParams{
					method:    "Infof",
					arguments: []any{"Success: GetNoteByUserIDAndNoteID with userID %v and noteID %v", up.userID, up.noteID},
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				mLogger.AssertCalled(t, wL.method, wL.arguments...)
			},
		},
		{
			name:         "GetNoteByUserIDandNoteID invalid userID",
			handler:      notesCtrl.GetNoteByUserIDAndNoteID,
			urlParams:    urlParams{userID: "-1", noteID: "1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantLogging: func(up urlParams) mockLoggingParams {
				return mockLoggingParams{
					method:    "Errorf",
					arguments: []any{"GetNoteByUserIDandNoteID: invalid userID %v", up.userID},
				}
			},
			wantBody: func(up urlParams) string { return "\n" },
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNoteByUserIDandNoteID")
				mLogger.AssertCalled(t, wL.method, wL.arguments...)
			},
		},
		{
			name:         "GetNoteByUserIDandNoteID invalid userID",
			handler:      notesCtrl.GetNoteByUserIDAndNoteID,
			urlParams:    urlParams{userID: "1", noteID: "-1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantBody:     func(up urlParams) string { return "\n" },
			wantLogging: func(up urlParams) mockLoggingParams {
				return mockLoggingParams{
					method:    "Errorf",
					arguments: []any{"GetNoteByUserIDandNoteID: invalid noteID %v", up.noteID},
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNoteByUserIDandNoteID")
				mLogger.AssertCalled(t, wL.method, wL.arguments...)
			},
		},
		{
			name:      "GetNoteByUserIDandNoteID DBError",
			handler:   notesCtrl.GetNoteByUserIDAndNoteID,
			urlParams: urlParams{userID: "1", noteID: "1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams {
				userID, noteID := mustConvUrlParamsToInt(t, up)
				return mockNotesStoreParams{
					method:          "GetNoteByUserIDAndNoteID",
					arguments:       []any{userID, noteID},
					returnArguments: []any{domain.Notes{}, errors.New("error notesStore.GetNoteByUserIDAndNoteID")},
				}
			},
			wantStatus: http.StatusNotFound,
			wantBody:   func(up urlParams) string { return "\n" },
			wantLogging: func(up urlParams) mockLoggingParams {
				msg := fmt.Sprintf("GetNoteByUserIDAndNoteID userID %v and noteID %v", up.userID, up.noteID)
				return mockLoggingParams{
					method: "Errorf",
					arguments: []any{
						"%s: %w",
						msg,
						errors.New("error notesStore.GetNoteByUserIDAndNoteID")},
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				mLogger.AssertCalled(t, wL.method, wL.arguments...)
			},
		},
	}

	for _, tc := range testCases {
		mNotesStore.Setup(tc.mockNSParams(tc.urlParams))
		mLogger.Setup(tc.wantLogging(tc.urlParams))
		req := setupRequest(t, "/users/{userID}/notes/{noteID}", tc.urlParams, &bytes.Reader{})
		rr := httptest.NewRecorder()
		tc.handler(rr, req)
		tc.assertions(t,
			rr,
			tc.wantStatus,
			tc.wantBody(tc.urlParams),
			tc.wantLogging(tc.urlParams),
			tc.mockNSParams(tc.urlParams))
	}
}

func TestGetNotesByUserID(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	mLogger := &mockLogger{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore, mLogger)

	testCases := []struct {
		name         string
		handler      http.HandlerFunc
		urlParams    urlParams
		mockNSParams func(urlP urlParams) mockNotesStoreParams
		wantStatus   int
		wantBody     func(urlParams) string
		wantLogging  func(urlParams) mockLoggingParams
		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams)
	}{
		{
			name:      "GetNotesByUserID success",
			handler:   notesCtrl.GetNotesByUserID,
			urlParams: urlParams{userID: "1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams {
				userID, _ := mustConvUrlParamsToInt(t, up)
				return mockNotesStoreParams{
					method:    "GetNotesByUserID",
					arguments: []any{userID},
					returnArguments: []any{domain.Notes{
						{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
						{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
					}, nil},
				}
			},
			wantStatus: http.StatusOK,
			wantBody: func(up urlParams) string {
				return mustEncode(t, domain.Notes{
					{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
					{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
				}).String()
			},
			wantLogging: func(up urlParams) mockLoggingParams {
				return mockLoggingParams{
					method:    "Infof",
					arguments: []any{"Success: GetNotesByUserID with userID %v", up.userID},
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				mLogger.AssertCalled(t, wL.method, wL.arguments...)
			},
		},
		{
			name:         "GetNotesByUserID invalid userID: not an int",
			handler:      notesCtrl.GetNotesByUserID,
			urlParams:    urlParams{userID: "bullshit"},
			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantLogging: func(up urlParams) mockLoggingParams {
				return mockLoggingParams{
					method:    "Errorf",
					arguments: []any{"GetNotesByUserID: invalid userID %v", up.userID},
				}
			},
			wantBody: func(up urlParams) string { return "\n" },
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNotesByUserID")
				mLogger.AssertCalled(t, wL.method, wL.arguments...)
			},
		},
		{
			name:         "GetNotesByUserID invalid userID: negative number",
			handler:      notesCtrl.GetNotesByUserID,
			urlParams:    urlParams{userID: "-1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantLogging: func(up urlParams) mockLoggingParams {
				return mockLoggingParams{
					method:    "Errorf",
					arguments: []any{"GetNotesByUserID: invalid userID %v", up.userID},
				}
			},
			wantBody: func(up urlParams) string { return "\n" },
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNotesByUserID")
				mLogger.AssertCalled(t, wL.method, wL.arguments...)
			},
		},
	}

	for _, tc := range testCases {
		mNotesStore.Setup(tc.mockNSParams(tc.urlParams))
		mLogger.Setup(tc.wantLogging(tc.urlParams))
		req := setupRequest(t, "/users/{userID}/notes", tc.urlParams, &bytes.Reader{})
		rr := httptest.NewRecorder()
		tc.handler(rr, req)
		tc.assertions(t,
			rr,
			tc.wantStatus,
			tc.wantBody(tc.urlParams),
			tc.wantLogging(tc.urlParams),
			tc.mockNSParams(tc.urlParams))
	}
}

func TestAdd(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	mLogger := &mockLogger{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore, mLogger)

	type NotePost struct {
		note string
	}

	testCases := []struct {
		name         string
		handler      http.HandlerFunc
		urlParams    urlParams
		body         io.Reader
		mockNSParams func(urlP urlParams) mockNotesStoreParams
		wantStatus   int
		wantBody     func(urlParams) string
		wantLogging  func(urlParams) mockLoggingParams
		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL mockLoggingParams, callAssertion mockNotesStoreParams)
	}{
		{
			name:       "Add Success",
			handler:    notesCtrl.Add,
			urlParams:  urlParams{userID: "1"},
			body:       mustEncode(t, NotePost{note: "Test note"}),
			wantStatus: http.StatusOK,
			mockNSParams: func(urlP urlParams) mockNotesStoreParams {

			},
		},
	}

	for _, tc := range testCases {
		mNotesStore.Setup(tc.mockNSParams(tc.urlParams))
		mLogger.Setup(tc.wantLogging(tc.urlParams))
		req := setupRequest(t, "/users/{userID}/notes", tc.urlParams, tc.body)
		rr := httptest.NewRecorder()
		tc.handler(rr, req)
		tc.assertions(t,
			rr,
			tc.wantStatus,
			tc.wantBody(tc.urlParams),
			tc.wantLogging(tc.urlParams),
			tc.mockNSParams(tc.urlParams))
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

func mustEncode(t *testing.T, a any) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(a); err != nil {
		t.Fatalf("encoding json: %v", err)
	}
	return buf
}

func setupRequest(t *testing.T, target string, up urlParams, body io.Reader) *http.Request {
	t.Helper()
	req := httptest.NewRequest("GET", target, body)
	return WithUrlParams(req, Params{
		"userID": up.userID,
		"noteID": up.noteID,
	})
}

func mustConvUrlParamsToInt(t *testing.T, up urlParams) (userID, noteID int) {
	t.Helper()
	userID, err := strconv.Atoi(up.userID)
	if err != nil {
		t.Fatal(err)
	}
	noteID, err = strconv.Atoi(up.userID)
	if err != nil {
		t.Fatal(err)
	}
	return userID, noteID
}
