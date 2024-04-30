package mid_test

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/Keisn1/note-taking-app/domain/web/auth"
	"github.com/Keisn1/note-taking-app/domain/web/mid"
	"github.com/Keisn1/note-taking-app/foundation/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuth struct {
	mock.Mock
}

type StubNoteRepo struct {
	notes map[uuid.UUID]note.Note
}

func (nr StubNoteRepo) Delete(noteID uuid.UUID) error                   { return nil }
func (nr StubNoteRepo) Create(n note.Note) error                        { return nil }
func (nr StubNoteRepo) Update(note note.Note) error                     { return nil }
func (nr StubNoteRepo) GetNoteByID(noteID uuid.UUID) (note.Note, error) { return nr.notes[noteID], nil }

func (nr StubNoteRepo) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) { return nil, nil }

func Test_Authorize(t *testing.T) {
	userID := uuid.New()
	noteID := uuid.New()
	snr := &StubNoteRepo{
		notes: map[uuid.UUID]note.Note{
			noteID: note.NewNote(noteID, note.Title{}, note.Content{}, userID),
		},
	}
	ns := note.NewNotesService(snr)

	midAuthorize := mid.AuthorizeNote(ns)
	handler := midAuthorize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Test Handler")) }))

	req := httptest.NewRequest(http.MethodGet, "/notImplemented", nil)
	req.SetPathValue("note_id", noteID.String())
	req = req.WithContext(context.WithValue(req.Context(), mid.UserIDKey, userID))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	req = httptest.NewRequest(http.MethodGet, "/notImplemented", nil)
	req.SetPathValue("note_id", noteID.String())
	falseUserID := uuid.New()
	req = req.WithContext(context.WithValue(req.Context(), mid.UserIDKey, falseUserID))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
	// testCases := []struct {
	// 	name      string
	// 	userID    uuid.UUID
	// 	noteID    uuid.UUID
	// 	assertion func(t *testing.T, err error)
	// }{}
	// for _, tc := range testCases {
	// 	t.Run(tc.name, func(t *testing.T) {
	// 		handler.ServeHTTP(rr, req)
	// 		tc.assertion(t, err)
	// 	})
	// }

}

func Test_Authenticate(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)

	key := common.MustGenerateRandomKey(32)
	jwtSvc, err := auth.NewJWTService(key)
	assert.NoError(t, err)
	a := auth.NewAuth(jwtSvc)

	midAuthenticate := mid.Authenticate(a)
	handler := midAuthenticate(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Test Handler")) }),
	)

	userID := uuid.New()
	testCases := []struct {
		name        string
		setupHeader func(*http.Request)
		assertions  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Test authentication success",
			setupHeader: func(req *http.Request) {
				tokenS, _ := jwtSvc.CreateToken(userID, time.Minute)
				req.Header.Set("Authorization", "Bearer "+tokenS)
			},
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Equal(t, "Test Handler", recorder.Body.String())
			},
		},
		{
			name:        "missing authentication header",
			setupHeader: func(req *http.Request) {},
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, recorder.Code)
				assert.Contains(t, recorder.Body.String(), "failed authentication")
				assert.Contains(t, logBuf.String(), "failed authentication")
			},
		},
		{
			name:        "invalid authentication header",
			setupHeader: func(req *http.Request) { req.Header.Set("Authorization", "invalid") },
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, recorder.Code)
				assert.Contains(t, recorder.Body.String(), "failed authentication")
				assert.Contains(t, logBuf.String(), "failed authentication")
			},
		},
		{
			name: "expired token",
			setupHeader: func(req *http.Request) {
				tokenS, _ := jwtSvc.CreateToken(userID, -1*time.Minute)
				req.Header.Set("Authorization", "Bearer "+tokenS)
			},
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, recorder.Code)
				assert.Contains(t, recorder.Body.String(), "failed authentication")
				assert.Contains(t, logBuf.String(), "failed authentication")
			},
		},
		{
			name:        "invalid token",
			setupHeader: func(req *http.Request) { req.Header.Set("Authorization", "Bearer invalid") },
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, recorder.Code)
				assert.Contains(t, recorder.Body.String(), "failed authentication")
				assert.Contains(t, logBuf.String(), "failed authentication")
			},
		},
	}

	for _, tc := range testCases {
		logBuf.Reset()
		req := httptest.NewRequest(http.MethodGet, "/auth", nil)
		tc.setupHeader(req)

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		tc.assertions(t, recorder)
	}
}
