package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"context"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthenticationMiddleware(t *testing.T) {
	// Initialize your JWT middleware and other necessary dependencies for testing

	// Create a new test server with the JWT middleware applied to the handler
	handler := JWTAuthenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest("GET", "/protected-route", nil)

	// Add a valid or invalid JWT token to the request headers for testing different scenarios
	ctx := context.Background()
	ctx = context.WithValue(ctx, JWTToken("token"), false)
	req = req.WithContext(ctx)

	// Make a request to the test server
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	// Assert the expected outcome based on the token validity
	assert.Equal(t, http.StatusForbidden, recorder.Code)
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
