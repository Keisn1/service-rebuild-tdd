package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestAuthentication(t *testing.T) {
	a := &Auth{}

	t.Run("Test Authentication happy path", func(t *testing.T) {
		userID := "123"
		current := time.Now()
		inOneHour := current.Add(1 * time.Hour)
		exp_time := jwt.NewNumericDate(inOneHour)

		issuer := os.Getenv("JWT_NOTES_ISSUER")
		claims := jwt.MapClaims{
			"sub": "123",
			"iss": issuer,
			"exp": exp_time,
		}

		secret := os.Getenv("JWT_SECRET_KEY")
		bearerToken := setupJwtTokenString(t, claims, secret)
		err := a.Authenticate(userID, bearerToken)

		assert.NoError(t, err)
	})

	t.Run("Test Authentication Failures", func(t *testing.T) {
		testBearerTokens := []string{
			"", "Bearer invalid length", "NoBearer asdf;lkj",
		}
		for _, bearerT := range testBearerTokens {
			err := a.Authenticate("", bearerT)
			ErrorContainss(t, err, "expected authorization header format: Bearer <token>", "authenticate:")
		}

		wrongMethodToken := getTokenEcdsa256(t)
		err := a.Authenticate("", "Bearer "+wrongMethodToken)
		ErrorContainss(t, err, "error parsing tokenString", "unexpected signing method: ES256", "authenticate:")

		invalidToken := "invalidToken"
		err = a.Authenticate("", "Bearer "+invalidToken)
		ErrorContainss(t, err, "error parsing tokenString", "authenticate:")

		secret := os.Getenv("JWT_SECRET_KEY")
		bearerToken := setupJwtTokenString(t, jwt.MapClaims{}, secret)
		err = a.Authenticate("", bearerToken)
		ErrorContainss(t, err, "no expiration date set", "authenticate:")

		oneMinuteAgo := jwt.NewNumericDate(time.Now().Add(-1 * time.Minute))
		claims := setupClaims(oneMinuteAgo, "", "")
		bearerToken = setupJwtTokenString(t, claims, secret)
		err = a.Authenticate("", bearerToken)
		ErrorContainss(t, err, "token is expired", "authenticate:")

		inOneHour := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
		claims = setupClaims(inOneHour, "false issuer", "")
		bearerToken = setupJwtTokenString(t, claims, secret)
		err = a.Authenticate("", bearerToken)
		ErrorContainss(t, err, "incorrect Issuer", "authenticate:")

		rightUserID, falseUserID := "123", "456"
		claims = setupClaims(inOneHour, os.Getenv("JWT_NOTES_ISSUER"), falseUserID)
		bearerToken = setupJwtTokenString(t, claims, secret)
		err = a.Authenticate(rightUserID, bearerToken)
		ErrorContainss(t, err, "user not enabled", "authenticate:")
	})
}

func ErrorContainss(t *testing.T, err error, containss ...string) {
	t.Helper()
	for _, contains := range containss {
		assert.ErrorContains(t, err, contains)
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

func setupJwtTokenString(t *testing.T, claims jwt.MapClaims, secret string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)
	bearerToken := "Bearer " + tokenS
	return bearerToken
}

func setupClaims(exp *jwt.NumericDate, iss, sub string) jwt.MapClaims {
	claims := jwt.MapClaims{
		"exp": exp,
		"iss": iss,
		"sub": sub,
	}
	return claims
}
