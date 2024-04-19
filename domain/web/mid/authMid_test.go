package mid_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestJWTAuthenticationMiddleware(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)

	key := common.MustGenerateRandomKey(32)
	jwtSvc, err := auth.NewJWTService(key)
	assert.NoError(t, err)
	a := auth.NewAuth(jwtSvc)

	authMid := mid.Authenticate(a)
	handler := authMid(http.HandlerFunc(
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
				assert.Equal(t, recorder.Body.String(), "Failed Authentication\n")
				assert.Contains(t, logBuf.String(), "Failed Authentication")
			},
		},
		{
			name:        "invalid authentication header",
			setupHeader: func(req *http.Request) { req.Header.Set("Authorization", "invalid") },
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, recorder.Code)
				assert.Equal(t, recorder.Body.String(), "Failed Authentication\n")
				assert.Contains(t, logBuf.String(), "Failed Authentication")
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
				assert.Equal(t, recorder.Body.String(), "Failed Authentication\n")
				assert.Contains(t, logBuf.String(), "Failed Authentication")
			},
		},
		{
			name:        "invalid token",
			setupHeader: func(req *http.Request) { req.Header.Set("Authorization", "Bearer invalid") },
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, recorder.Code)
				assert.Equal(t, recorder.Body.String(), "Failed Authentication\n")
				assert.Contains(t, logBuf.String(), "Failed Authentication")
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
