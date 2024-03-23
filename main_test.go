package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"

	"crypto/ecdsa"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestAuthentication(t *testing.T) {
	a := &Auth{}
	t.Run("Test false authorization header format", func(t *testing.T) {
		testBearerTokens := []string{
			"", "Bearer invalid length", "NoBearer asdf;lkj",
		}
		for _, bearerT := range testBearerTokens {
			_, err := a.getTokenString(bearerT)
			assert.Contains(t, err.Error(), "expected authorization header format: Bearer <token>")
		}
	})

	t.Run("Test wrong signing method", func(t *testing.T) {
		wrongMethodToken := getTokenEcdsa256(t)
		_, err := a.parseTokenString(wrongMethodToken)
		assert.Contains(t, err.Error(), "unexpected signing method: ES256")
	})

	t.Run("Test invalid token", func(t *testing.T) {
		invalidToken := "invalidToken"
		_, err := a.parseTokenString(invalidToken)
		assert.Error(t, err)
	})

	t.Run("Test that user is NOT enabled", func(t *testing.T) {
		userID := "123"
		claims := jwt.MapClaims{
			"sub": "456",
		}
		err := a.isUserEnabled(userID, claims)
		assert.ErrorContains(t, err, "user not enabled")
	})

	t.Run("Test that user IS enabled", func(t *testing.T) {
		userID := "123"
		claims := jwt.MapClaims{
			"sub": "123",
		}
		err := a.isUserEnabled(userID, claims)
		assert.NoError(t, err)
	})

	t.Run("Test Authentication pipeline happy path", func(t *testing.T) {
		userID := "123"
		wantClaims := jwt.MapClaims{
			"sub": "123",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, wantClaims)
		tokenS, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
		assert.NoError(t, err)
		bearerToken := "Bearer " + tokenS
		gotClaims, err := a.Authenticate(userID, bearerToken)
		assert.NoError(t, err)
		assert.Equal(t, gotClaims, wantClaims)
	})
}

func TestJWTAuthenticationMiddleware(t *testing.T) {
	handler := JWTAuthenticationMiddleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Test Handler"))
		}),
	)

	t.Run("Test false authorization header format", func(t *testing.T) {
		testReqs := []*http.Request{
			newEmptyGetRequest(t),
			addAuthorizationJWT(t, "invalid length", newEmptyGetRequest(t)),
			addFalseAuthorizationHeader(t, "", newEmptyGetRequest(t)),
		}
		for _, req := range testReqs {
			var logBuf bytes.Buffer
			log.SetOutput(&logBuf)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusForbidden, recorder.Code)
			assert.Contains(t, recorder.Body.String(), "Failed Authorization")
			assert.Contains(t, logBuf.String(), "expected authorization header format: Bearer <token>")
		}
	})

	t.Run("Test invalid signing method", func(t *testing.T) {
		tString := getTokenEcdsa256(t) // wrong signing method
		req := newEmptyGetRequest(t)
		req = addAuthorizationJWT(t, tString, req)

		var logBuf bytes.Buffer
		log.SetOutput(&logBuf)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusForbidden, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Failed Authorization")
		assert.Contains(t, logBuf.String(), "unexpected signing method")
	})

	t.Run("Test invalid token", func(t *testing.T) {
		tString := "InvalidToken"
		req := newEmptyGetRequest(t)
		req = addAuthorizationJWT(t, tString, req)

		var logBuf bytes.Buffer
		log.SetOutput(&logBuf)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusForbidden, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Failed Authorization")
		assert.Contains(t, logBuf.String(), "Token invalid")
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

func addAuthorizationJWT(t *testing.T, tokenS string, req *http.Request) *http.Request {
	req.Header.Add("Authorization", "Bearer "+tokenS)
	return req
}

func addFalseAuthorizationHeader(t *testing.T, tokenS string, req *http.Request) *http.Request {
	req.Header.Add("Authorization", "False "+tokenS)
	return req
}
