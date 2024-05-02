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
	"github.com/Keisn1/note-taking-app/foundation"
	"github.com/Keisn1/note-taking-app/foundation/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Authorize(t *testing.T) {
	userID := uuid.New()
	noteID := uuid.New()
	n := note.NewNote(noteID, note.Title{}, note.Content{}, userID)
	sns := &StubNoteService{notes: map[uuid.UUID]note.Note{noteID: n}}
	midAuthorize := mid.AuthorizeNote(sns)

	t.Run("Authorize success, retreived note set in context of request", func(t *testing.T) {
		wantNote := n
		wantResponse := "test handler"
		handler := midAuthorize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(wantResponse))

			gotNote := mid.GetNote(r.Context())
			assert.Equal(t, wantNote, gotNote)
		}))

		req := httptest.NewRequest(http.MethodGet, "/notImplemented", nil)
		req.SetPathValue("note_id", noteID.String())

		req = req.WithContext(context.WithValue(req.Context(), foundation.UserIDKey, userID))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, wantResponse, rr.Body.String())
	})

	t.Run("Authorize failure, wrong user id", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/notImplemented", nil)
		req.SetPathValue("note_id", noteID.String())
		falseUserID := uuid.New()
		req = req.WithContext(context.WithValue(req.Context(), foundation.UserIDKey, falseUserID))

		handler := midAuthorize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Test Handler"))
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("Authorize failure, note not present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/notImplemented", nil)
		req.SetPathValue("note_id", uuid.New().String())
		userID := uuid.New()
		req = req.WithContext(context.WithValue(req.Context(), foundation.UserIDKey, userID))

		handler := midAuthorize(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Test Handler"))
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})
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

	t.Run("Test claims set on context after success", func(t *testing.T) {
		key := common.MustGenerateRandomKey(32)
		jwtSvc, err := auth.NewJWTService(key)
		assert.NoError(t, err)
		a := auth.NewAuth(jwtSvc)

		midAuthenticate := mid.Authenticate(a)

		wantUserID := uuid.New()
		tokenS, err := jwtSvc.CreateToken(wantUserID, time.Minute)
		assert.NoError(t, err)

		wantClaims, err := a.Authenticate("Bearer " + tokenS)
		assert.NoError(t, err)

		handler := midAuthenticate(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				gotUserID := mid.GetUserID(r.Context())
				assert.Equal(t, wantUserID, gotUserID)

				gotClaims := mid.GetClaims(r.Context())
				assert.Equal(t, wantClaims, gotClaims)
			}),
		)

		req := httptest.NewRequest(http.MethodGet, "/auth", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+tokenS)

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
	})

}
