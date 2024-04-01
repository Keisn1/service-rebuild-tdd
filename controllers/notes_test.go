package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"fmt"

	"errors"

	ctrls "github.com/Keisn1/note-taking-app/controllers"
	"github.com/Keisn1/note-taking-app/domain"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type urlParams struct {
	userID string
	noteID string
}

func TestGetAllNotes(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
	logBuf := &bytes.Buffer{}
	log.SetOutput(logBuf)

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
		wantLogging  []string
		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams)
	}{
		{
			name:         "GetAllNotes Success",
			handler:      notesCtrl.GetAllNotes,
			mockNSParams: mockNotesStoreParams{method: "GetAllNotes", returnArguments: []any{notes, nil}},
			wantStatus:   http.StatusOK,
			wantBody:     mustEncode(t, notes).String(),
			wantLogging:  []string{"Success: GetAllNotes"},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, rr.Body.String(), wantBody)
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)

			},
		},
		{
			name:    "GetAllNotes Error DB",
			handler: notesCtrl.GetAllNotes,
			mockNSParams: mockNotesStoreParams{
				method:          "GetAllNotes",
				returnArguments: []any{domain.Notes{}, errors.New("error notesStore.GetAllNotes")},
			},
			wantBody:    "\n",
			wantStatus:  http.StatusInternalServerError,
			wantLogging: []string{"ERROR", "GetAllNotes: DBError", "error notesStore.GetAllNotes"},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, rr.Body.String(), wantBody)
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
	}

	for _, tc := range testCases {
		logBuf.Reset()
		mNotesStore.Setup(tc.mockNSParams)
		req := setupRequest(t, "GET", "/users/notes", urlParams{}, &bytes.Buffer{})
		rr := httptest.NewRecorder()
		notesCtrl.GetAllNotes(rr, req)
		tc.assertions(t, rr, tc.wantStatus, tc.wantBody, tc.wantLogging, tc.mockNSParams)
	}
}

func TestGetNoteByUserIDandNoteID(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
	logBuf := &bytes.Buffer{}
	log.SetOutput(logBuf)

	testCases := []struct {
		name         string
		handler      http.HandlerFunc
		urlParams    urlParams
		mockNSParams func(urlP urlParams) mockNotesStoreParams
		wantStatus   int
		wantBody     func(urlParams) string
		wantLogging  func(urlParams) []string
		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams)
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
			wantLogging: func(up urlParams) []string {
				return []string{
					fmt.Sprintf("Success: GetNoteByUserIDAndNoteID: userID %v noteID %v", up.userID, up.noteID),
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:         "GetNoteByUserIDandNoteID invalid userID",
			handler:      notesCtrl.GetNoteByUserIDAndNoteID,
			urlParams:    urlParams{userID: "-1", noteID: "1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantLogging: func(up urlParams) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("GetNoteByUserIDandNoteID: invalid userID %v", up.userID),
				}
			},
			wantBody: func(up urlParams) string { return "\n" },
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNoteByUserIDandNoteID")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:         "GetNoteByUserIDandNoteID invalid noteID",
			handler:      notesCtrl.GetNoteByUserIDAndNoteID,
			urlParams:    urlParams{userID: "1", noteID: "-1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantLogging: func(up urlParams) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("GetNoteByUserIDandNoteID: invalid noteID %v", up.noteID),
				}
			},
			wantBody: func(up urlParams) string { return "\n" },
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNoteByUserIDandNoteID")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
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
			wantLogging: func(up urlParams) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("GetNoteByUserIDAndNoteID: userID %v noteID %v", up.userID, up.noteID),
					"error notesStore.GetNoteByUserIDAndNoteID",
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
	}

	for _, tc := range testCases {
		logBuf.Reset()
		mNotesStore.Setup(tc.mockNSParams(tc.urlParams))
		req := setupRequest(t, "GET", "/users/{userID}/notes/{noteID}", tc.urlParams, &bytes.Buffer{})
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
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
	logBuf := &bytes.Buffer{}
	log.SetOutput(logBuf)

	testCases := []struct {
		name         string
		handler      http.HandlerFunc
		urlParams    urlParams
		mockNSParams func(urlP urlParams) mockNotesStoreParams
		wantStatus   int
		wantBody     func(urlParams) string
		wantLogging  func(urlParams) []string
		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams)
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
			wantLogging: func(up urlParams) []string {
				return []string{
					"INFO",
					fmt.Sprintf("Success: GetNotesByUserID: userID %v", up.userID),
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:         "GetNotesByUserID invalid userID: not an int",
			handler:      notesCtrl.GetNotesByUserID,
			urlParams:    urlParams{userID: "bullshit"},
			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantLogging: func(up urlParams) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("GetNotesByUserID: invalid userID %v", up.userID),
				}
			},
			wantBody: func(up urlParams) string { return "\n" },
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNotesByUserID")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:         "GetNotesByUserID invalid userID: negative number",
			handler:      notesCtrl.GetNotesByUserID,
			urlParams:    urlParams{userID: "-1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantLogging: func(up urlParams) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("GetNotesByUserID: invalid userID %v", up.userID),
				}
			},
			wantBody: func(up urlParams) string { return "\n" },
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNotesByUserID")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:      "GetNotesByUserID DBError",
			handler:   notesCtrl.GetNotesByUserID,
			urlParams: urlParams{userID: "1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams {
				userID, _ := mustConvUrlParamsToInt(t, up)
				return mockNotesStoreParams{
					method:          "GetNotesByUserID",
					arguments:       []any{userID},
					returnArguments: []any{domain.Notes{}, errors.New("error notesStore.GetNotesByUserID")},
				}
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   func(up urlParams) string { return "\n" },
			wantLogging: func(up urlParams) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("GetNotesByUserID: userID %v", up.userID),
					"error notesStore.GetNotesByUserID",
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
	}

	for _, tc := range testCases {
		logBuf.Reset()
		mNotesStore.Setup(tc.mockNSParams(tc.urlParams))
		req := setupRequest(t, "GET", "/users/{userID}/notes", tc.urlParams, &bytes.Buffer{})
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

func TestAddNote(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
	logBuf := &bytes.Buffer{}
	log.SetOutput(logBuf)

	testCases := []struct {
		name         string
		handler      http.HandlerFunc
		urlParams    urlParams
		body         domain.NotePost
		reqBody      func(body domain.NotePost) *bytes.Buffer
		mockNSParams func(urlP urlParams, body domain.NotePost) mockNotesStoreParams
		wantStatus   int
		wantBody     string
		wantLogging  func(urlParams, domain.NotePost) []string
		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams)
	}{
		{
			name:      "Add Success",
			handler:   notesCtrl.Add,
			urlParams: urlParams{userID: "1"},
			body:      domain.NotePost{Note: "Test note"},
			reqBody:   func(body domain.NotePost) *bytes.Buffer { return mustEncode(t, body) },
			mockNSParams: func(up urlParams, body domain.NotePost) mockNotesStoreParams {
				userID, _ := mustConvUrlParamsToInt(t, up)
				return mockNotesStoreParams{
					method: "AddNote", arguments: []any{userID, body.Note}, returnArguments: []any{nil},
				}
			},
			wantStatus: http.StatusAccepted,
			wantLogging: func(up urlParams, body domain.NotePost) []string {
				return []string{
					"INFO",
					fmt.Sprintf("Success: Add: userID %v note %v", up.userID, body),
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:         "AddNote invalid userID",
			handler:      notesCtrl.Add,
			urlParams:    urlParams{userID: "bullshit"},
			reqBody:      func(body domain.NotePost) *bytes.Buffer { return mustEncode(t, nil) },
			mockNSParams: func(up urlParams, body domain.NotePost) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantBody:     "\n",
			wantLogging: func(up urlParams, body domain.NotePost) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("Add: invalid userID %v", up.userID),
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "AddNote")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:         "AddNote invalid userID negative number",
			handler:      notesCtrl.Add,
			urlParams:    urlParams{userID: "-1"},
			reqBody:      func(body domain.NotePost) *bytes.Buffer { return mustEncode(t, nil) },
			mockNSParams: func(up urlParams, body domain.NotePost) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantBody:     "\n",
			wantLogging: func(up urlParams, body domain.NotePost) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("Add: invalid userID %v", up.userID),
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "AddNote")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:      "Add DBError",
			handler:   notesCtrl.Add,
			urlParams: urlParams{userID: "1"},
			reqBody:   func(body domain.NotePost) *bytes.Buffer { return mustEncode(t, nil) },
			mockNSParams: func(up urlParams, body domain.NotePost) mockNotesStoreParams {
				userID, _ := mustConvUrlParamsToInt(t, up)
				return mockNotesStoreParams{
					method:          "AddNote",
					arguments:       []any{userID, body.Note},
					returnArguments: []any{errors.New("error notesStore.AddNote")},
				}
			},
			wantStatus: http.StatusConflict,
			wantBody:   "\n",
			wantLogging: func(up urlParams, body domain.NotePost) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("Add: userID %v body %v", up.userID, body),
					"error notesStore.AddNote",
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "AddNote")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:         "Add with invalid body",
			handler:      notesCtrl.Add,
			urlParams:    urlParams{userID: "1"},
			reqBody:      func(body domain.NotePost) *bytes.Buffer { return bytes.NewBuffer([]byte("invalid body")) },
			mockNSParams: func(up urlParams, body domain.NotePost) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantBody:     "\n",
			wantLogging: func(up urlParams, np domain.NotePost) []string {
				return []string{"ERROR", "Add: invalid body"}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, callAssertion.method, callAssertion.arguments...)
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
	}

	for _, tc := range testCases {
		logBuf.Reset()
		mNotesStore.Setup(tc.mockNSParams(tc.urlParams, tc.body))
		req := setupRequest(t, "POST", "/users/{userID}/notes", tc.urlParams, tc.reqBody(tc.body))
		rr := httptest.NewRecorder()
		tc.handler(rr, req)
		tc.assertions(
			t,
			rr,
			tc.wantStatus,
			tc.wantBody,
			tc.wantLogging(tc.urlParams, tc.body),
			tc.mockNSParams(tc.urlParams, tc.body),
		)
	}
}

func TestDelete(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
	logBuf := &bytes.Buffer{}
	log.SetOutput(logBuf)

	testCases := []struct {
		name         string
		handler      http.HandlerFunc
		urlParams    urlParams
		mockNSParams func(urlP urlParams) mockNotesStoreParams
		wantStatus   int
		wantBody     string
		wantLogging  func(urlParams) []string
		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams)
	}{
		{
			name:      "Success Deletion",
			handler:   notesCtrl.Delete,
			urlParams: urlParams{userID: "1", noteID: "1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams {
				userID, noteID := mustConvUrlParamsToInt(t, up)
				return mockNotesStoreParams{
					method:          "Delete",
					arguments:       []any{userID, noteID},
					returnArguments: []any{nil},
				}
			},
			wantStatus: http.StatusNoContent,
			wantLogging: func(up urlParams) []string {
				return []string{
					"INFO",
					fmt.Sprintf("Success: Delete: userID %v noteID %v", up.userID, up.noteID),
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:         "Delete invalid userID",
			handler:      notesCtrl.Delete,
			urlParams:    urlParams{userID: "-1", noteID: "1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantLogging: func(up urlParams) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("Delete: invalid userID %v", up.userID),
				}
			},
			wantBody: "\n",
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "Delete")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:         "Delete invalid noteID",
			handler:      notesCtrl.Delete,
			urlParams:    urlParams{userID: "1", noteID: "-1"},
			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:   http.StatusBadRequest,
			wantLogging: func(up urlParams) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("Delete: invalid noteID %v", up.noteID),
				}
			},
			wantBody: "\n",
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "Delete")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
	}

	for _, tc := range testCases {
		logBuf.Reset()
		mNotesStore.Setup(tc.mockNSParams(tc.urlParams))
		req := setupRequest(t, "DELETE", "/users/{userID}/notes", tc.urlParams, &bytes.Buffer{})
		rr := httptest.NewRecorder()
		tc.handler(rr, req)
		tc.assertions(
			t,
			rr,
			tc.wantStatus,
			tc.wantBody,
			tc.wantLogging(tc.urlParams),
			tc.mockNSParams(tc.urlParams),
		)
	}
}

// func TestNotes2(t *testing.T) {
// 	notes := domain.Notes{
// 		{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
// 		{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
// 		{NoteID: 3, UserID: 2, Note: "Note 1 user 2"},
// 		{NoteID: 4, UserID: 2, Note: "Note 2 user 2"},
// 	}

// 	notesStore := NewStubNotesStore(notes)
// 	logger := NewStubLogger()
// 	notesCtrl := ctrls.NewNotesCtrlr(notesStore, logger)

// 	// t.Run("Deletion fail", func(t *testing.T) {
// 	// 	logger.Reset()

// 	// 	userID, noteID := 50, 50
// 	// 	request, err := http.NewRequest(http.MethodDelete, "", nil)
// 	// 	assertNoError(t, err)
// 	// 	request = WithUrlParams(request, Params{
// 	// 		"userID": strconv.Itoa(userID),
// 	// 		"noteID": strconv.Itoa(noteID),
// 	// 	})

// 	// 	response := httptest.NewRecorder()
// 	// 	notesCtrl.Delete(response, request)

// 	// 	assertStatusCode(t, response.Result().StatusCode, http.StatusInternalServerError)
// 	// 	assertLoggingCalls(t, logger.errorfCall, []string{"Delete DBerror:"})
// 	// })

// 	t.Run("Edit a Note", func(t *testing.T) {
// 		logger.Reset()

// 		userID, noteID, note := 1, 1, "New note text"
// 		putRequest := newPutRequestWithNoteAndUrlParams(t, "New note text", Params{
// 			"userID": strconv.Itoa(userID),
// 			"noteID": strconv.Itoa(noteID),
// 		})
// 		response := httptest.NewRecorder()
// 		notesCtrl.Edit(response, putRequest)

// 		wantEditCalls := []EditCall{{userID: userID, noteID: noteID, note: note}}
// 		assertStatusCode(t, response.Result().StatusCode, http.StatusOK)
// 		assertSlicesAnyAreEqual(t, notesStore.editNoteCalls, wantEditCalls)
// 		assertLoggingCalls(t, logger.infofCalls, []string{
// 			"Success: Edit: userID 1 noteID 1 note New note text",
// 		})
// 	})
// }

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

func setupRequest(t *testing.T, method, target string, up urlParams, body *bytes.Buffer) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, target, body)
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
