package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	a := &Auth{}
	secret := os.Getenv("JWT_SECRET_KEY")
	issuer := os.Getenv("JWT_NOTES_ISSUER")

	testCases := []struct {
		name        string
		userID      string
		bearerToken func() string
		assertion   func(t *testing.T, err error)
	}{
		{
			name:        "Empty Bearer",
			bearerToken: func() string { return "" },
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: expected authorization header format: Bearer <token>")
			},
		},
		{
			name:        "Wrong format length",
			bearerToken: func() string { return "Bearer invalid length" },
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: expected authorization header format: Bearer <token>")
			},
		},
		{
			name:        "Wrong format Prefix",
			bearerToken: func() string { return "NoBearer asdf;lkj" },
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: expected authorization header format: Bearer <token>")
			},
		},
		{
			name:        "Wrong method",
			bearerToken: func() string { return getBearerTokenEcdsa256(t) },
			assertion: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "authenticate: error parsing tokenString")
			},
		},
		{
			name:        "Invalid Token",
			bearerToken: func() string { return "Bearer invalidToken" },
			assertion: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "authenticate: error parsing tokenString")
			},
		},
		{
			name: "False secret",
			bearerToken: func() string {
				return setupJwtTokenString(t, jwt.MapClaims{}, "falseSecret")
			},
			assertion: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "authenticate: error parsing tokenString")
			},
		},
		{
			name: "Expired Token",
			bearerToken: func() string {
				oneMinuteAgo := jwt.NewNumericDate(time.Now().Add(-1 * time.Minute))
				claims := setupClaims(oneMinuteAgo, "", "")
				return setupJwtTokenString(t, claims, secret)
			},
			assertion: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "authenticate: error parsing tokenString")
			},
		},
		{
			name: "No expiration date set",
			bearerToken: func() string {
				return setupJwtTokenString(t, jwt.MapClaims{}, secret)
			},
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: no expiration date set")
			},
		},
		{
			name: "False issuer",
			bearerToken: func() string {
				inOneHour := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
				claims := setupClaims(inOneHour, "false issuer", "")
				return setupJwtTokenString(t, claims, secret)
			},
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: incorrect issuer")
			},
		},
		{
			name:   "userID NOT equal subject in jwt (456)",
			userID: "123",
			bearerToken: func() string {
				inOneHour := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
				claims := setupClaims(inOneHour, issuer, "456")
				return setupJwtTokenString(t, claims, secret)
			},
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: invalid subject")
			},
		},
		{
			name:   "userID in endpoint unequal userID in jwt",
			userID: "123",
			bearerToken: func() string {
				inOneHour := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
				claims := setupClaims(inOneHour, issuer, "123")
				return setupJwtTokenString(t, claims, secret)
			},
			assertion: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := a.Authenticate(tc.userID, tc.bearerToken())
			tc.assertion(t, err)
		})
	}
}

func getBearerTokenEcdsa256(t *testing.T) (tokenString string) {
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
	return "Bearer " + tokenString
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
