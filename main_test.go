package main

import (
	"log"
	"reflect"
	"testing"
)

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
