package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"context"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthenticationMiddleware(t *testing.T) {
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
	handler := JWTAuthenticationMiddleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Test Handler"))
		}),
	)

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
}

func TestAuthentication(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	gotClaims, err := ValidateToken(tokenString)
	assertNoError(t, err)
	wantClaims := map[string]string{
		"iss": "note-taking-app",
	}
	if !reflect.DeepEqual(gotClaims, wantClaims) {
		t.Errorf("got = %v; want %v", gotClaims, wantClaims)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		log.Fatal(err)
	}
}
