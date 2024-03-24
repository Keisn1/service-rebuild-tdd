package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bytes"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthentication(t *testing.T) {
	a := &Auth{}
	t.Run("Test Authentication Errors", func(t *testing.T) {
		testBearerTokens := []string{
			"", "Bearer invalid length", "NoBearer asdf;lkj",
		}
		for _, bearerT := range testBearerTokens {
			_, err := a.Authenticate("", bearerT)
			assert.ErrorContains(t, err, "expected authorization header format: Bearer <token>")
			assert.ErrorContains(t, err, "authenticate:")
		}

		wrongMethodToken := getTokenEcdsa256(t)
		_, err := a.Authenticate("", "Bearer "+wrongMethodToken)
		assert.ErrorContains(t, err, "unexpected signing method: ES256")
		assert.ErrorContains(t, err, "error parsing tokenString")
		assert.ErrorContains(t, err, "authenticate:")

		invalidToken := "invalidToken"
		_, err = a.Authenticate("", "Bearer "+invalidToken)
		assert.ErrorContains(t, err, "error parsing tokenString")
		assert.ErrorContains(t, err, "authenticate:")

		userID := "123"
		claims := jwt.MapClaims{
			"sub": "456",
		}
		bearerToken := setupJwtTokenString(t, claims)
		_, err = a.Authenticate(userID, bearerToken)
		assert.ErrorContains(t, err, "user not enabled")
		assert.ErrorContains(t, err, "authenticate:")
	})

	t.Run("Test Authentication pipeline happy path", func(t *testing.T) {
		userID := "123"
		wantClaims := jwt.MapClaims{
			"sub": "123",
		}

		bearerToken := setupJwtTokenString(t, wantClaims)
		gotClaims, err := a.Authenticate(userID, bearerToken)

		assert.NoError(t, err)
		assert.Equal(t, gotClaims, wantClaims)
	})
}

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
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Test Handler"))
		}),
	)

	t.Run("Test auth authenticate is called", func(t *testing.T) {
		mockAuth.On("Authenticate", "123", "Bearer valid token").Return(jwt.MapClaims{}, nil)

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		req = WithUrlParam(req, "userID", "123")

		req.Header.Set("Authorization", "Bearer valid token")
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		mockAuth.AssertCalled(t, "Authenticate", "123", "Bearer valid token")
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("Test authentication failure", func(t *testing.T) {
		var claims jwt.MapClaims
		mockAuth.On("Authenticate", "123", "Bearer INVALID token").Return(
			claims, errors.New("error in authenticate"),
		)

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer INVALID token")
		req = WithUrlParam(req, "userID", "123")

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

func getTokenEcdsa256(t *testing.T) (tokenString string) {
	t.Helper()
	var (
		key   *ecdsa.PrivateKey
		token *jwt.Token
	)
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	token = jwt.New(jwt.SigningMethodES256)
	tokenString, err = token.SignedString(key)
	assert.NoError(t, err)
	return tokenString
}

func newEmptyGetRequest(t *testing.T) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	return req
}

func addAuthorizationJWT(t *testing.T, tokenS string, req *http.Request) *http.Request {
	req.Header.Add("Authorization", "Bearer "+tokenS)
	return req
}

func setupJwtTokenString(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	assert.NoError(t, err)
	bearerToken := "Bearer " + tokenS
	return bearerToken
}

func WithUrlParam(r *http.Request, key, value string) *http.Request {
	chiCtx := chi.NewRouteContext()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add(key, value)
	return r
}
