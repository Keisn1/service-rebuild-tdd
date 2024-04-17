package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"os"
	"testing"

	"github.com/Keisn1/note-taking-app/foundation/jwt"
	jwtLib "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	secret := os.Getenv("JWT_SECRET_KEY")
	a := &Auth{jwt: jwt.NewJWT(secret)}

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
			name:        "Failing token verification",
			bearerToken: func() string { return getBearerTokenEcdsa256(t) },
			assertion: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "authenticate: verify: ")
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
		token *jwtLib.Token
	)
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	token = jwtLib.New(jwtLib.SigningMethodES256)
	tokenString, err = token.SignedString(key)
	assert.NoError(t, err)
	return "Bearer " + tokenString
}

func setupJwtTokenString(t *testing.T, claims jwtLib.MapClaims, secret string) string {
	t.Helper()
	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
	tokenS, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)
	bearerToken := "Bearer " + tokenS
	return bearerToken
}

func setupClaims(exp *jwtLib.NumericDate, iss, sub string) jwtLib.MapClaims {
	claims := jwtLib.MapClaims{
		"exp": exp,
		"iss": iss,
		"sub": sub,
	}
	return claims
}
