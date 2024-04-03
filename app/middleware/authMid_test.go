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
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)

	mockAuth := new(MockAuth)
	jwtMidHandler := NewJwtMidHandler(mockAuth)
	handler := jwtMidHandler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Test Handler")) }),
	)

	testCases := []struct {
		name          string
		setupMockAuth func()
		setupRequest  func() *http.Request
		assertions    func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Test authentication success",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/auth", nil)
				req = WithUrlParam(req, "userID", "123")
				req.Header.Set("Authorization", "Bearer valid token")
				return req
			},
			setupMockAuth: func() {
				mockAuth.On("Authenticate", "123", "Bearer valid token").Return(jwt.MapClaims{}, nil)
			},
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Equal(t, "Test Handler", recorder.Body.String())
				mockAuth.AssertCalled(t, "Authenticate", "123", "Bearer valid token")
			},
		},
		{
			name: "Test authentication Failure",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/auth", nil)
				req = WithUrlParam(req, "userID", "123")
				req.Header.Set("Authorization", "Bearer INVALID token")
				return req
			},
			setupMockAuth: func() {
				var claims jwt.MapClaims
				mockAuth.On("Authenticate", "123", "Bearer INVALID token").Return(
					claims, errors.New("error in authenticate"),
				)
			},
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, recorder.Code)
				assert.Contains(t, recorder.Body.String(), "Failed Authentication")
				assert.Contains(t, logBuf.String(),
					"Failed Authentication userID \"123\" bearerToken \"Bearer INVALID token\"")
			},
		},
	}

	for _, tc := range testCases {
		logBuf.Reset()
		tc.setupMockAuth()
		req := tc.setupRequest()

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		tc.assertions(t, recorder)
	}
}

func WithUrlParam(r *http.Request, key, value string) *http.Request {
	chiCtx := chi.NewRouteContext()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add(key, value)
	return r
}
