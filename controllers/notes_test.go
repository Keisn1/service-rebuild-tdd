package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"errors"
	ctrls "github.com/Keisn1/note-taking-app/controllers"
	"github.com/Keisn1/note-taking-app/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddNote(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
	logBuf := &bytes.Buffer{}
	log.SetOutput(logBuf)

	type testCase struct {
		name        string
		userID      any
		body        any
		mNSP        func(userID any, body any) mockNotesStoreParams
		wantStatus  int
		wantBody    string
		wantLogging func(userID any, body any) []string
		assertions  func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams)
	}

	setupRequest := func(t *testing.T, method, target string, userID any, body *bytes.Buffer) *http.Request {
		t.Helper()
		req := httptest.NewRequest(method, target, body)
		ctx := context.WithValue(req.Context(), ctrls.UserIDKey, userID)
		req = req.WithContext(ctx)
		return req
	}

	testCases := []testCase{
		{
			name:   "Add Success",
			userID: uuid.New(),
			body:   domain.NotePost{Note: "Test note"},
			mNSP: func(userID any, body any) mockNotesStoreParams {
				return mockNotesStoreParams{method: "AddNote", arguments: []any{userID, body}, returnArguments: []any{nil}}
			},
			wantStatus: http.StatusAccepted,
			wantBody:   "",
			wantLogging: func(userID any, body any) []string {
				return []string{
					"INFO",
					fmt.Sprintf("Success: Add: userID %v body %v", userID, body)}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertCalled(t, mNSP.method, mNSP.arguments...)
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:        "Add: Invalid userID",
			userID:      "invalid userID",
			body:        domain.NotePost{Note: "Test note"},
			mNSP:        func(userID any, body any) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus:  http.StatusBadRequest,
			wantBody:    "\n",
			wantLogging: func(userID any, body any) []string { return []string{"ERROR", "Add: invalid userID"} },
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "AddNote")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:   "Add DBError",
			userID: uuid.New(),
			body:   domain.NotePost{Note: "Test note"},
			mNSP: func(userID any, body any) mockNotesStoreParams {
				return mockNotesStoreParams{
					method:          "AddNote",
					arguments:       []any{userID, body},
					returnArguments: []any{errors.New("error notesStore.AddNote")},
				}
			},
			wantStatus: http.StatusConflict,
			wantBody:   "\n",
			wantLogging: func(userID any, body any) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("Add: userID %v body %v", userID, body),
					"error notesStore.AddNote",
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "AddNote")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:       "Add with invalid body",
			userID:     uuid.New(),
			body:       "invalid body",
			mNSP:       func(userID any, body any) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus: http.StatusBadRequest,
			wantBody:   "\n",
			wantLogging: func(userID any, body any) []string {
				return []string{"ERROR", "Add: invalid body"}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "AddNote")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
	}

	for _, tc := range testCases {
		logBuf.Reset()
		mNotesStore.Setup(tc.mNSP(tc.userID, tc.body))
		req := setupRequest(t, "POST", "/users/notes", tc.userID, mustEncode(t, tc.body))
		rr := httptest.NewRecorder()
		notesCtrl.Add(rr, req)
		tc.assertions(t, rr, tc.wantStatus, tc.wantBody, tc.wantLogging(tc.userID, tc.body), tc.mNSP(tc.userID, tc.body))
	}
}

// func TestGetAllNotes(t *testing.T) {
// 	mNotesStore := &mockNotesStore{}
// 	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
// 	logBuf := &bytes.Buffer{}
// 	log.SetOutput(logBuf)

// 	notes := domain.Notes{
// 		{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
// 		{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
// 		{NoteID: 3, UserID: 2, Note: "Note 1 user 2"},
// 		{NoteID: 4, UserID: 2, Note: "Note 2 user 2"},
// 	}
// 	testCases := []struct {
// 		name         string
// 		handler      http.HandlerFunc
// 		mockNSParams mockNotesStoreParams
// 		wantStatus   int
// 		wantBody     string
// 		wantLogging  []string
// 		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams)
// 	}{
// 		{
// 			name:         "GetAllNotes Success",
// 			handler:      notesCtrl.GetAllNotes,
// 			mockNSParams: mockNotesStoreParams{method: "GetAllNotes", returnArguments: []any{notes, nil}},
// 			wantStatus:   http.StatusOK,
// 			wantBody:     mustEncode(t, notes).String(),
// 			wantLogging:  []string{"Success: GetAllNotes"},
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, rr.Body.String(), wantBody)
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)

// 			},
// 		},
// 		{
// 			name:    "GetAllNotes Error DB",
// 			handler: notesCtrl.GetAllNotes,
// 			mockNSParams: mockNotesStoreParams{
// 				method:          "GetAllNotes",
// 				returnArguments: []any{domain.Notes{}, errors.New("error notesStore.GetAllNotes")},
// 			},
// 			wantBody:    "\n",
// 			wantStatus:  http.StatusInternalServerError,
// 			wantLogging: []string{"ERROR", "GetAllNotes: DBError", "error notesStore.GetAllNotes"},
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, rr.Body.String(), wantBody)
// 				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		logBuf.Reset()
// 		mNotesStore.Setup(tc.mockNSParams)
// 		req := setupRequest(t, "GET", "/users/notes", urlParams{}, &bytes.Buffer{})
// 		rr := httptest.NewRecorder()
// 		notesCtrl.GetAllNotes(rr, req)
// 		tc.assertions(t, rr, tc.wantStatus, tc.wantBody, tc.wantLogging, tc.mockNSParams)
// 	}
// }

// func TestGetNoteByUserIDandNoteID(t *testing.T) {
// 	mNotesStore := &mockNotesStore{}
// 	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
// 	logBuf := &bytes.Buffer{}
// 	log.SetOutput(logBuf)

// 	testCases := []struct {
// 		name         string
// 		handler      http.HandlerFunc
// 		urlParams    urlParams
// 		body         domain.NotePost
// 		mockNSParams func(urlP urlParams) mockNotesStoreParams
// 		wantStatus   int
// 		wantBody     func(urlParams) string
// 		wantLogging  func(urlParams) []string
// 		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams)
// 	}{
// 		{
// 			name:      "GetNoteByUserIDandNoteID success",
// 			handler:   notesCtrl.GetNoteByUserIDAndNoteID,
// 			urlParams: urlParams{userID: "1", noteID: "1"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams {
// 				userID, noteID := mustConvUrlParamsToInt(t, up)
// 				note := fmt.Sprintf("Note %d user %d", userID, noteID)
// 				return mockNotesStoreParams{
// 					method:    "GetNoteByUserIDAndNoteID",
// 					arguments: []any{userID, noteID},
// 					returnArguments: []any{
// 						domain.Notes{{UserID: userID, NoteID: noteID, Note: note}},
// 						nil,
// 					},
// 				}
// 			},
// 			wantStatus: http.StatusOK,
// 			wantBody: func(up urlParams) string {
// 				userID, noteID := mustConvUrlParamsToInt(t, up)
// 				note := fmt.Sprintf("Note %d user %d", userID, noteID)
// 				return mustEncode(t, domain.Notes{{UserID: userID, NoteID: noteID, Note: note}}).String()
// 			},
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					fmt.Sprintf("Success: GetNoteByUserIDAndNoteID: userID %v noteID %v", up.userID, up.noteID),
// 				}
// 			},
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 		{
// 			name:         "GetNoteByUserIDandNoteID invalid userID",
// 			handler:      notesCtrl.GetNoteByUserIDAndNoteID,
// 			urlParams:    urlParams{userID: "-1", noteID: "1"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
// 			wantStatus:   http.StatusBadRequest,
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					"ERROR",
// 					fmt.Sprintf("GetNoteByUserIDandNoteID: invalid userID %v", up.userID),
// 				}
// 			},
// 			wantBody: func(up urlParams) string { return "\n" },
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertNotCalled(t, "GetNoteByUserIDandNoteID")
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 		{
// 			name:         "GetNoteByUserIDandNoteID invalid noteID",
// 			handler:      notesCtrl.GetNoteByUserIDAndNoteID,
// 			urlParams:    urlParams{userID: "1", noteID: "-1"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
// 			wantStatus:   http.StatusBadRequest,
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					"ERROR",
// 					fmt.Sprintf("GetNoteByUserIDandNoteID: invalid noteID %v", up.noteID),
// 				}
// 			},
// 			wantBody: func(up urlParams) string { return "\n" },
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertNotCalled(t, "GetNoteByUserIDandNoteID")
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 		{
// 			name:      "GetNoteByUserIDandNoteID DBError",
// 			handler:   notesCtrl.GetNoteByUserIDAndNoteID,
// 			urlParams: urlParams{userID: "1", noteID: "1"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams {
// 				userID, noteID := mustConvUrlParamsToInt(t, up)
// 				return mockNotesStoreParams{
// 					method:          "GetNoteByUserIDAndNoteID",
// 					arguments:       []any{userID, noteID},
// 					returnArguments: []any{domain.Notes{}, errors.New("error notesStore.GetNoteByUserIDAndNoteID")},
// 				}
// 			},
// 			wantStatus: http.StatusInternalServerError,
// 			wantBody:   func(up urlParams) string { return "\n" },
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					"ERROR",
// 					fmt.Sprintf("GetNoteByUserIDAndNoteID: userID %v noteID %v", up.userID, up.noteID),
// 					"error notesStore.GetNoteByUserIDAndNoteID",
// 				}
// 			},
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 		{
// 			name:      "GetNoteByUserIDandNoteID Not found if empty return from DB",
// 			handler:   notesCtrl.GetNoteByUserIDAndNoteID,
// 			urlParams: urlParams{userID: "1", noteID: "1"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams {
// 				userID, noteID := mustConvUrlParamsToInt(t, up)
// 				return mockNotesStoreParams{
// 					method:          "GetNoteByUserIDAndNoteID",
// 					arguments:       []any{userID, noteID},
// 					returnArguments: []any{domain.Notes{}, nil},
// 				}
// 			},
// 			wantStatus: http.StatusNotFound,
// 			wantBody:   func(up urlParams) string { return "\n" },
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					"ERROR",
// 					fmt.Sprintf("GetNoteByUserIDAndNoteID: userID %v noteID %v", up.userID, up.noteID),
// 					"Not Found",
// 				}
// 			},
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		logBuf.Reset()
// 		mNotesStore.Setup(tc.mockNSParams(tc.urlParams))
// 		req := setupRequest(t, "GET", "/users/{userID}/notes/{noteID}", tc.urlParams, &bytes.Buffer{})
// 		rr := httptest.NewRecorder()
// 		tc.handler(rr, req)
// 		tc.assertions(t,
// 			rr,
// 			tc.wantStatus,
// 			tc.wantBody(tc.urlParams),
// 			tc.wantLogging(tc.urlParams),
// 			tc.mockNSParams(tc.urlParams))
// 	}
// }

// func TestGetNotesByUserID(t *testing.T) {
// 	mNotesStore := &mockNotesStore{}
// 	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
// 	logBuf := &bytes.Buffer{}
// 	log.SetOutput(logBuf)

// 	testCases := []struct {
// 		name         string
// 		handler      http.HandlerFunc
// 		urlParams    urlParams
// 		mockNSParams func(urlP urlParams) mockNotesStoreParams
// 		wantStatus   int
// 		wantBody     func(urlParams) string
// 		wantLogging  func(urlParams) []string
// 		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams)
// 	}{
// 		{
// 			name:      "GetNotesByUserID success",
// 			handler:   notesCtrl.GetNotesByUserID,
// 			urlParams: urlParams{userID: "1"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams {
// 				userID, _ := mustConvUrlParamsToInt(t, up)
// 				return mockNotesStoreParams{
// 					method:    "GetNotesByUserID",
// 					arguments: []any{userID},
// 					returnArguments: []any{domain.Notes{
// 						{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
// 						{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
// 					}, nil},
// 				}
// 			},
// 			wantStatus: http.StatusOK,
// 			wantBody: func(up urlParams) string {
// 				return mustEncode(t, domain.Notes{
// 					{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
// 					{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
// 				}).String()
// 			},
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					"INFO",
// 					fmt.Sprintf("Success: GetNotesByUserID: userID %v", up.userID),
// 				}
// 			},
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 		{
// 			name:         "GetNotesByUserID invalid userID: not an int",
// 			handler:      notesCtrl.GetNotesByUserID,
// 			urlParams:    urlParams{userID: "bullshit"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
// 			wantStatus:   http.StatusBadRequest,
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					"ERROR",
// 					fmt.Sprintf("GetNotesByUserID: invalid userID %v", up.userID),
// 				}
// 			},
// 			wantBody: func(up urlParams) string { return "\n" },
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertNotCalled(t, "GetNotesByUserID")
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 		{
// 			name:         "GetNotesByUserID invalid userID: negative number",
// 			handler:      notesCtrl.GetNotesByUserID,
// 			urlParams:    urlParams{userID: "-1"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
// 			wantStatus:   http.StatusBadRequest,
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					"ERROR",
// 					fmt.Sprintf("GetNotesByUserID: invalid userID %v", up.userID),
// 				}
// 			},
// 			wantBody: func(up urlParams) string { return "\n" },
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertNotCalled(t, "GetNotesByUserID")
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 		{
// 			name:      "GetNotesByUserID DBError",
// 			handler:   notesCtrl.GetNotesByUserID,
// 			urlParams: urlParams{userID: "1"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams {
// 				userID, _ := mustConvUrlParamsToInt(t, up)
// 				return mockNotesStoreParams{
// 					method:          "GetNotesByUserID",
// 					arguments:       []any{userID},
// 					returnArguments: []any{domain.Notes{}, errors.New("error notesStore.GetNotesByUserID")},
// 				}
// 			},
// 			wantStatus: http.StatusInternalServerError,
// 			wantBody:   func(up urlParams) string { return "\n" },
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					"ERROR",
// 					fmt.Sprintf("GetNotesByUserID: userID %v", up.userID),
// 					"error notesStore.GetNotesByUserID",
// 				}
// 			},
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		logBuf.Reset()
// 		mNotesStore.Setup(tc.mockNSParams(tc.urlParams))
// 		req := setupRequest(t, "GET", "/users/{userID}/notes", tc.urlParams, &bytes.Buffer{})
// 		rr := httptest.NewRecorder()
// 		tc.handler(rr, req)
// 		tc.assertions(t,
// 			rr,
// 			tc.wantStatus,
// 			tc.wantBody(tc.urlParams),
// 			tc.wantLogging(tc.urlParams),
// 			tc.mockNSParams(tc.urlParams))
// 	}
// }

// func TestDelete(t *testing.T) {
// 	mNotesStore := &mockNotesStore{}
// 	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
// 	logBuf := &bytes.Buffer{}
// 	log.SetOutput(logBuf)

// 	testCases := []struct {
// 		name         string
// 		handler      http.HandlerFunc
// 		urlParams    urlParams
// 		mockNSParams func(urlP urlParams) mockNotesStoreParams
// 		wantStatus   int
// 		wantBody     string
// 		wantLogging  func(urlParams) []string
// 		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams)
// 	}{
// 		{
// 			name:      "Success Deletion",
// 			handler:   notesCtrl.Delete,
// 			urlParams: urlParams{userID: "1", noteID: "1"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams {
// 				userID, noteID := mustConvUrlParamsToInt(t, up)
// 				return mockNotesStoreParams{
// 					method:          "Delete",
// 					arguments:       []any{userID, noteID},
// 					returnArguments: []any{nil},
// 				}
// 			},
// 			wantStatus: http.StatusNoContent,
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					"INFO",
// 					fmt.Sprintf("Success: Delete: userID %v noteID %v", up.userID, up.noteID),
// 				}
// 			},
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertCalled(t, callAssertion.method, callAssertion.arguments...)
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 		{
// 			name:         "Delete invalid userID",
// 			handler:      notesCtrl.Delete,
// 			urlParams:    urlParams{userID: "-1", noteID: "1"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
// 			wantStatus:   http.StatusBadRequest,
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					"ERROR",
// 					fmt.Sprintf("Delete: invalid userID %v", up.userID),
// 				}
// 			},
// 			wantBody: "\n",
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertNotCalled(t, "Delete")
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 		{
// 			name:         "Delete invalid noteID",
// 			handler:      notesCtrl.Delete,
// 			urlParams:    urlParams{userID: "1", noteID: "-1"},
// 			mockNSParams: func(up urlParams) mockNotesStoreParams { return mockNotesStoreParams{} },
// 			wantStatus:   http.StatusBadRequest,
// 			wantLogging: func(up urlParams) []string {
// 				return []string{
// 					"ERROR",
// 					fmt.Sprintf("Delete: invalid noteID %v", up.noteID),
// 				}
// 			},
// 			wantBody: "\n",
// 			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams) {
// 				assert.Equal(t, wantStatus, rr.Code)
// 				assert.Equal(t, wantBody, rr.Body.String())
// 				mNotesStore.AssertNotCalled(t, "Delete")
// 				for _, logMsg := range wL {
// 					assert.Contains(t, logBuf.String(), logMsg)
// 				}
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		logBuf.Reset()
// 		mNotesStore.Setup(tc.mockNSParams(tc.urlParams))
// 		req := setupRequest(t, "DELETE", "/users/{userID}/notes", tc.urlParams, &bytes.Buffer{})
// 		rr := httptest.NewRecorder()
// 		tc.handler(rr, req)
// 		tc.assertions(
// 			t,
// 			rr,
// 			tc.wantStatus,
// 			tc.wantBody,
// 			tc.wantLogging(tc.urlParams),
// 			tc.mockNSParams(tc.urlParams),
// 		)
// 	}
// }

// WithUrlParam returns a pointer to a request object with the given URL params
// added to a new chi.Context object.
type Params map[string]string

func WithUrlParams(r *http.Request, params Params) *http.Request {
	chiCtx := chi.NewRouteContext()
	req := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	for key, value := range params {
		chiCtx.URLParams.Add(key, value)
	}
	return req
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

// func mustConvUrlParamsToInt(t *testing.T, up urlParams) (userID, noteID int) {
// 	t.Helper()
// 	userID, err := strconv.Atoi(up.userID)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	noteID, err = strconv.Atoi(up.userID)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	return userID, noteID
// }
