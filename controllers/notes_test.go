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

	testCases := []testCase{
		{
			name:   "Add Success",
			userID: uuid.New(),
			body:   domain.NotePost{Note: "Test note"},
			mNSP: func(userID any, body any) mockNotesStoreParams {
				np := body.(domain.NotePost)
				return mockNotesStoreParams{method: "AddNote", arguments: []any{userID, np.Note}, returnArguments: []any{nil}}
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
			wantBody:    fmt.Sprintln(""),
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
				np := body.(domain.NotePost)
				return mockNotesStoreParams{
					method:          "AddNote",
					arguments:       []any{userID, np.Note},
					returnArguments: []any{errors.New("error notesStore.AddNote")},
				}
			},
			wantStatus: http.StatusConflict,
			wantBody:   fmt.Sprintln(""),
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
			wantBody:   fmt.Sprintln(""),
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
		req := setupRequest(t, "POST", "/users/notes", tc.userID, "", mustEncode(t, tc.body))
		rr := httptest.NewRecorder()
		notesCtrl.Add(rr, req)
		tc.assertions(t, rr, tc.wantStatus, tc.wantBody, tc.wantLogging(tc.userID, tc.body), tc.mNSP(tc.userID, tc.body))
	}
}

func TestGetAllNotes(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
	logBuf := &bytes.Buffer{}
	log.SetOutput(logBuf)

	type testCase struct {
		name        string
		mNSP        mockNotesStoreParams
		wantStatus  int
		wantBody    string
		wantLogging []string
		assertions  func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams)
	}

	testNotes := domain.Notes{
		{NoteID: 1, UserID: uuid.New(), Note: "Note 1 user 1"},
		{NoteID: 2, UserID: uuid.New(), Note: "Note 2 user 1"},
		{NoteID: 3, UserID: uuid.New(), Note: "Note 1 user 2"},
		{NoteID: 4, UserID: uuid.New(), Note: "Note 2 user 2"},
	}
	testCases := []testCase{

		{
			name:        "GetAllNotes Success",
			mNSP:        mockNotesStoreParams{method: "GetAllNotes", returnArguments: []any{testNotes, nil}},
			wantStatus:  http.StatusOK,
			wantBody:    mustEncode(t, testNotes).String(),
			wantLogging: []string{"Success: GetAllNotes"},
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
			name: "GetAllNotes Error DB",
			mNSP: mockNotesStoreParams{
				method:          "GetAllNotes",
				returnArguments: []any{domain.Notes{}, errors.New("error notesStore.GetAllNotes")},
			},
			wantBody:    fmt.Sprintln(""),
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
		mNotesStore.Setup(tc.mNSP)
		req := setupRequest(t, "GET", "/notes", nil, "", &bytes.Buffer{})
		rr := httptest.NewRecorder()
		notesCtrl.GetAllNotes(rr, req)
		tc.assertions(t, rr, tc.wantStatus, tc.wantBody, tc.wantLogging, tc.mNSP)
	}
}

func TestGetNotesByUserID(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
	logBuf := &bytes.Buffer{}
	log.SetOutput(logBuf)

	type testCase struct {
		name        string
		userID      any
		mNSP        func(userID any) mockNotesStoreParams
		wantStatus  int
		wantBody    func(userID any) string
		wantLogging func(userID any) []string
		assertions  func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams)
	}

	testCases := []testCase{
		{
			name:   "GetNotesByUserID success",
			userID: uuid.New(),
			mNSP: func(userID any) mockNotesStoreParams {
				uid := userID.(uuid.UUID)
				return mockNotesStoreParams{
					method:    "GetNotesByUserID",
					arguments: []any{userID},
					returnArguments: []any{domain.Notes{
						{NoteID: 1, UserID: uid, Note: "Note 1 user 1"},
						{NoteID: 2, UserID: uid, Note: "Note 2 user 1"},
					}, nil},
				}
			},
			wantStatus: http.StatusOK,
			wantBody: func(userID any) string {
				uid := userID.(uuid.UUID)
				return mustEncode(t, domain.Notes{
					{NoteID: 1, UserID: uid, Note: "Note 1 user 1"},
					{NoteID: 2, UserID: uid, Note: "Note 2 user 1"},
				}).String()
			},
			wantLogging: func(userID any) []string {
				return []string{
					"INFO",
					fmt.Sprintf("Success: GetNotesByUserID: userID %v", userID),
				}
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
			name:       "GetNotesByUserID invalid userID",
			userID:     123,
			mNSP:       func(userID any) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus: http.StatusBadRequest,
			wantBody:   func(userID any) string { return fmt.Sprintln("") },
			wantLogging: func(userID any) []string {
				return []string{
					"ERROR",
					"GetNotesByUserID: invalid userID",
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNotesByUserID")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:   "GetNotesByUserID DBError",
			userID: uuid.New(),
			mNSP: func(userID any) mockNotesStoreParams {
				return mockNotesStoreParams{
					method:          "GetNotesByUserID",
					arguments:       []any{userID},
					returnArguments: []any{domain.Notes{}, errors.New("error notesStore.GetNotesByUserID")},
				}
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   func(userID any) string { return fmt.Sprintln("") },
			wantLogging: func(userID any) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("GetNotesByUserID: userID %v", userID),
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
		mNotesStore.Setup(tc.mNSP(tc.userID))
		req := setupRequest(t, "GET", "/users/notes", tc.userID, "", &bytes.Buffer{})
		rr := httptest.NewRecorder()
		notesCtrl.GetNotesByUserID(rr, req)
		tc.assertions(t, rr, tc.wantStatus, tc.wantBody(tc.userID), tc.wantLogging(tc.userID), tc.mNSP(tc.userID))
	}
}

func TestGetNoteByUserIDandNoteID(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
	logBuf := &bytes.Buffer{}
	log.SetOutput(logBuf)

	type testCase struct {
		name        string
		userID      any
		noteID      string
		mNSP        func(userID any, noteID string) mockNotesStoreParams
		wantStatus  int
		wantBody    func(userID any, noteID string) string
		wantLogging func(userID any, noteID string) []string
		assertions  func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams)
	}

	testCases := []testCase{
		{
			name:   "GetNoteByUserIDandNoteID success",
			userID: uuid.New(),
			noteID: "1",
			mNSP: func(userID any, noteID string) mockNotesStoreParams {
				uid := userID.(uuid.UUID)
				nid := mustAtoi(t, noteID)
				return mockNotesStoreParams{
					method:    "GetNoteByUserIDAndNoteID",
					arguments: []any{uid, nid},
					returnArguments: []any{
						domain.Notes{{UserID: uid, NoteID: nid, Note: "test note"}},
						nil,
					},
				}
			},
			wantStatus: http.StatusOK,
			wantBody: func(userID any, noteID string) string {
				uid := userID.(uuid.UUID)
				nid := mustAtoi(t, noteID)
				return mustEncode(t, domain.Notes{{UserID: uid, NoteID: nid, Note: "test note"}}).String()
			},
			wantLogging: func(userID any, noteID string) []string {
				return []string{fmt.Sprintf("Success: GetNoteByUserIDAndNoteID: userID %v noteID %v", userID, noteID)}
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
			name:       "GetNoteByUserIDandNoteID invalid userID",
			userID:     "invalid userID",
			noteID:     "1",
			mNSP:       func(userID any, noteID string) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus: http.StatusBadRequest,
			wantBody:   func(userID any, noteID string) string { return fmt.Sprintln("") },
			wantLogging: func(userID any, noteID string) []string {
				return []string{
					"ERROR",
					"GetNoteByUserIDandNoteID: invalid userID",
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNotesByUserID")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:       "GetNoteByUserIDandNoteID invalid noteID",
			userID:     uuid.New(),
			noteID:     "-1",
			mNSP:       func(userID any, noteID string) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus: http.StatusBadRequest,
			wantBody:   func(userID any, noteID string) string { return fmt.Sprintln("") },
			wantLogging: func(userID any, noteID string) []string {
				return []string{
					"ERROR",
					"GetNoteByUserIDandNoteID: invalid noteID",
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNotesByUserID")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:       "GetNoteByUserIDandNoteID invalid noteID",
			userID:     uuid.New(),
			noteID:     "invalid noteID",
			mNSP:       func(userID any, noteID string) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus: http.StatusBadRequest,
			wantBody:   func(userID any, noteID string) string { return fmt.Sprintln("") },
			wantLogging: func(userID any, noteID string) []string {
				return []string{
					"ERROR",
					"GetNoteByUserIDandNoteID: invalid noteID",
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "GetNotesByUserID")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},

		{
			name:   "GetNoteByUserIDandNoteID DBError",
			userID: uuid.New(),
			noteID: "1",
			mNSP: func(userID any, noteID string) mockNotesStoreParams {
				uid := userID.(uuid.UUID)
				nid := mustAtoi(t, noteID)
				return mockNotesStoreParams{
					method:    "GetNoteByUserIDAndNoteID",
					arguments: []any{uid, nid},
					returnArguments: []any{
						domain.Notes{},
						errors.New("error notesStore.GetNoteByUserIDAndNoteID"),
					},
				}
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   func(userID any, noteID string) string { return fmt.Sprintln("") },
			wantLogging: func(userID any, noteID string) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("GetNoteByUserIDAndNoteID: userID %v noteID %v", userID, noteID),
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
		{
			name:   "GetNoteByUserIDandNoteID Not found if empty return from DB",
			userID: uuid.New(),
			noteID: "1",
			mNSP: func(userID any, noteID string) mockNotesStoreParams {
				uid := userID.(uuid.UUID)
				nid := mustAtoi(t, noteID)
				return mockNotesStoreParams{
					method:          "GetNoteByUserIDAndNoteID",
					arguments:       []any{uid, nid},
					returnArguments: []any{domain.Notes{}, nil},
				}
			},
			wantStatus: http.StatusNotFound,
			wantBody:   func(userID any, noteID string) string { return fmt.Sprintln("") },
			wantLogging: func(userID any, noteID string) []string {
				return []string{
					"ERROR",
					fmt.Sprintf("GetNoteByUserIDAndNoteID: userID %v noteID %v", userID, noteID),
					"Not Found",
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
		mNotesStore.Setup(tc.mNSP(tc.userID, tc.noteID))

		req := setupRequest(t, "GET", "/users/notes/{noteID}", tc.userID, tc.noteID, &bytes.Buffer{})
		rr := httptest.NewRecorder()
		notesCtrl.GetNoteByUserIDAndNoteID(rr, req)

		tc.assertions(t, rr, tc.wantStatus, tc.wantBody(tc.userID, tc.noteID), tc.wantLogging(tc.userID, tc.noteID), tc.mNSP(tc.userID, tc.noteID))
	}
}

func TestDelete(t *testing.T) {
	mNotesStore := &mockNotesStore{}
	notesCtrl := ctrls.NewNotesCtrlr(mNotesStore)
	logBuf := &bytes.Buffer{}
	log.SetOutput(logBuf)

	type testCase struct {
		name        string
		userID      any
		noteID      string
		mNSP        func(userID any, noteID string) mockNotesStoreParams
		wantStatus  int
		wantBody    string
		wantLogging func(userID any, noteID string) []string
		assertions  func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams)
	}

	testCases := []testCase{
		{
			name:   "Success Deletion",
			userID: uuid.New(),
			noteID: "1",
			mNSP: func(userID any, noteID string) mockNotesStoreParams {
				uid := userID.(uuid.UUID)
				nid := mustAtoi(t, noteID)
				return mockNotesStoreParams{
					method:          "Delete",
					arguments:       []any{uid, nid},
					returnArguments: []any{nil},
				}
			},
			wantStatus: http.StatusNoContent,
			wantBody:   "",
			wantLogging: func(userID any, noteID string) []string {
				return []string{
					"INFO",
					fmt.Sprintf("Success: Delete: userID %v noteID %v", userID, noteID),
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
			name:       "Delete invalid userID",
			userID:     "invalid userID",
			noteID:     "1",
			mNSP:       func(userID any, noteID string) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus: http.StatusBadRequest,
			wantBody:   fmt.Sprintln(""),
			wantLogging: func(userID any, noteID string) []string {
				return []string{
					"ERROR",
					"Delete: invalid userID",
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "Delete")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:       "Delete invalid noteID",
			userID:     uuid.New(),
			noteID:     "-1",
			mNSP:       func(userID any, noteID string) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus: http.StatusBadRequest,
			wantBody:   fmt.Sprintln(""),
			wantLogging: func(userID any, noteID string) []string {
				return []string{
					"ERROR",
					"Delete: invalid noteID",
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams) {
				assert.Equal(t, wantStatus, rr.Code)
				assert.Equal(t, wantBody, rr.Body.String())
				mNotesStore.AssertNotCalled(t, "Delete")
				for _, logMsg := range wL {
					assert.Contains(t, logBuf.String(), logMsg)
				}
			},
		},
		{
			name:       "Delete invalid noteID",
			userID:     uuid.New(),
			noteID:     "invalid noteID",
			mNSP:       func(userID any, noteID string) mockNotesStoreParams { return mockNotesStoreParams{} },
			wantStatus: http.StatusBadRequest,
			wantBody:   fmt.Sprintln(""),
			wantLogging: func(userID any, noteID string) []string {
				return []string{
					"ERROR",
					"Delete: invalid noteID",
				}
			},
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, mNSP mockNotesStoreParams) {
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
		mNotesStore.Setup(tc.mNSP(tc.userID, tc.noteID))

		req := setupRequest(t, "DELETE", "/users/notes/{noteID}", tc.userID, tc.noteID, &bytes.Buffer{})
		rr := httptest.NewRecorder()
		notesCtrl.Delete(rr, req)

		tc.assertions(t, rr, tc.wantStatus, tc.wantBody, tc.wantLogging(tc.userID, tc.noteID), tc.mNSP(tc.userID, tc.noteID))
	}
}

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

func mustEncode(t *testing.T, a any) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(a); err != nil {
		t.Fatalf("encoding json: %v", err)
	}
	return buf
}

func mustAtoi(t *testing.T, s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		t.Fatal()
	}
	return i
}

func setupRequest(t *testing.T, method, target string, userID any, noteID string, body *bytes.Buffer) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, target, body)
	ctx := context.WithValue(req.Context(), ctrls.UserIDKey, userID)
	req = req.WithContext(ctx)
	return WithUrlParams(req, Params{
		"noteID": noteID,
	})
}
