package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"os"

	"bytes"

	"crypto/ecdsa"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthenticationMiddleware(t *testing.T) {
	handler := JWTAuthenticationMiddleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Test Handler"))
		}),
	)

	t.Run("Test invalid invalid signing method", func(t *testing.T) {
		tString := getTokenEcdsa256(t) // wrong signing method
		req := newEmptyGetRequest(t)
		req = addJwtTokenToContext(t, tString, req)

		var buf bytes.Buffer
		log.SetOutput(&buf)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusForbidden, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Invalid Authorization")
		assert.Contains(t, buf.String(), "unexpected signing method")
	})

	t.Run("Test invalid token", func(t *testing.T) {
		tString := "An Invalid string"
		req := newEmptyGetRequest(t)
		req = addJwtTokenToContext(t, tString, req)

		var buf bytes.Buffer
		log.SetOutput(&buf)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusForbidden, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Invalid Authorization")
		assert.Contains(t, buf.String(), "Invalid Authorization")
	})

	t.Run("Test token validation", func(t *testing.T) {
		// Initialize your JWT middleware and other necessary dependencies for testing
		secretKey := os.Getenv("JWT_SECRET_KEY")
		invalidTokenString := "An Invalid string"
		validTokenString, err := jwt.New(jwt.SigningMethodHS256).SignedString([]byte(secretKey))
		assertNoError(t, err)
		testCases := []struct {
			tokenString string
			statusCode  int
			wantBody    string
		}{
			{invalidTokenString, http.StatusForbidden, "No valid JWTToken"},
			{validTokenString, http.StatusOK, "Test Handler"},
		}

		// Create a new test server with the JWT middleware applied to the handler

		for _, tc := range testCases {
			req := httptest.NewRequest("GET", "/protected-route", nil)

			// Add a valid or invalid JWT token to the request headers for testing different scenarios
			req = req.WithContext(context.WithValue(context.Background(), JWTToken("token"), tc.tokenString))

			// Make a request to the test server
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, tc.statusCode, recorder.Code)
		}
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		log.Fatal(err)
	}
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

func addJwtTokenToContext(t *testing.T, tString string, req *http.Request) *http.Request {
	return req.WithContext(context.WithValue(context.Background(), JWTToken("token"), tString))
}
