package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"os"

	"bytes"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthenticationMiddleware(t *testing.T) {
	handler := JWTAuthenticationMiddleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Test Handler"))
		}),
	)

	t.Run("Test invalid token", func(t *testing.T) {
		invalidTokenString := "An Invalid string"

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assertNoError(t, err)
		req = req.WithContext(context.WithValue(context.Background(), JWTToken("token"), invalidTokenString))

		var buf bytes.Buffer
		log.SetOutput(&buf)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusForbidden, recorder.Code)
		assert.Equal(t, "Invalid JWT", buf.String())

		// test token valid
		// test claims valid
		// test claim userID
		// test userID equal url parameter userID
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
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		}
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		log.Fatal(err)
	}
}
