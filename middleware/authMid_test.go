package authMid

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuth struct {
	mock.Mock
}

func (ma *MockAuth) Authenticate(userID, bearerToken string) (jwt.Claims, error) {
	args := ma.Called(userID, bearerToken)
	return args.Get(0).(jwt.Claims), args.Error(1)
}

func TestJWTAuthenticationMiddleware(t *testing.T) {
	mockAuth := new(MockAuth)
	jwtMidHandler := NewJwtMidHandler(mockAuth)
	handler := jwtMidHandler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Test Handler")) }),
	)

	testCases := []struct {
		name          string
		setupMockAuth func()
		setupRequest  func(req *http.Request) *http.Request
	}{
		{
			name: "Test authentication success",
			setupRequest: func(req *http.Request) *http.Request {
				req = WithUrlParam(req, "userID", "123")
				req.Header.Set("Authorization", "Bearer valid token")
				return req
			},
			setupMockAuth: func() {
				mockAuth.On("Authenticate", "123", "Bearer valid token").Return(jwt.MapClaims{}, nil)
			},
		},
	}

	for _, tc := range testCases {
		tc.setupMockAuth()

		req := httptest.NewRequest(http.MethodGet, "/auth", nil)
		req = tc.setupRequest(req)

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		mockAuth.AssertCalled(t, "Authenticate", "123", "Bearer valid token")
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "Test Handler", recorder.Body.String())
	}

	t.Run("Test authentication success", func(t *testing.T) {
		var claims jwt.MapClaims
		mockAuth.On("Authenticate", "123", "Bearer INVALID token").Return(
			claims, errors.New("error in authenticate"),
		)

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		req = WithUrlParam(req, "userID", "123")
		req.Header.Set("Authorization", "Bearer INVALID token")

		var logBuf bytes.Buffer
		log.SetOutput(&logBuf)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		mockAuth.AssertCalled(t, "Authenticate", "123", "Bearer INVALID token")
		assert.Equal(t, http.StatusForbidden, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Failed Authentication")
		assert.Contains(t, logBuf.String(), "Failed Authentication")
		assert.Contains(t, logBuf.String(), "error in authenticate")
	})
}

func WithUrlParam(r *http.Request, key, value string) *http.Request {
	chiCtx := chi.NewRouteContext()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add(key, value)
	return r
}
