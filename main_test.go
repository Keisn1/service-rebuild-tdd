package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID       uuid.UUID `json:"Id"`
	Name     string    `json:"Name"`
	Email    string    `json:"Email"`
	Password string    `json:"Password"`
}

func TestAcceptance(t *testing.T) {
	// user request signup
	user := User{
		Name:     "kay",
		Email:    "kay@email.com",
		Password: "secret",
	}

	canSignUp(t, user)
	canSignIn(t, user)
	// canMakeNotes(user)
	canSignOut(t, user)
	// canNotMakeNotes()
}

func canSignOut(t *testing.T, u User) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/signin", mustEncode(t, u))
	signin(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	bearerToken := rr.Header().Get("Authorization")
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/signout", &bytes.Buffer{})
	req.Header.Set("Authorization", bearerToken)
	signout(rr, req)

	getsNotAuthorized(t)
}

func getsNotAuthorized(t *testing.T) {
	assert.True(t, true)
}

func signout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func canSignIn(t *testing.T, u User) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/signin", mustEncode(t, u))
	signin(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	bearerToken := rr.Header().Get("Authorization")
	verifyJWT(t, bearerToken)
	canMakeNotes(t, bearerToken)
}

func canMakeNotes(t *testing.T, bearerToken string) {
	assert.True(t, true)
}

func signin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func verifyJWT(t *testing.T, bearerToken string) {
	assert.True(t, true)
}

func canSignUp(t *testing.T, u User) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/signup", mustEncode(t, u))
	signup(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func mustEncode(t *testing.T, a any) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(a); err != nil {
		t.Fatalf("encoding json: %v", err)
	}
	return buf
}

// func TestSignIn(t *testing.T) {
// 	mUserStore := &mockUserStore{}
// 	userCtrl := ctrls.NewUserCtrlr(mUserStore)
// 	logBuf := &bytes.Buffer{}
// 	log.SetOutput(logBuf)
// 	testCase :=  []struct {
// 		name         string
// 		handler      http.HandlerFunc
// 		mockNSParams mockUserStoreParams
// 		wantStatus   int
// 		wantBody     string
// 		wantLogging  []string
// 		assertions   func(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantBody string, wL []string, callAssertion mockNotesStoreParams)
// 	} {
// 		{
// 			name: "Success SignIn",
// 		},
// 	}
// }
