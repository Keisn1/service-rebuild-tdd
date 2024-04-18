package jwtSvc_test

import (
	"testing"
	"time"

	"crypto/rand"

	"fmt"
	"github.com/Keisn1/note-taking-app/foundation/jwtSvc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	key := mustGenerateRandomKey(24)
	_, err := jwtSvc.NewJWT(key)
	assert.EqualError(t, err, "key minLength 32")

	key = mustGenerateRandomKey(32)
	jwtS, err := jwtSvc.NewJWT(key)
	assert.NoError(t, err)

	userID := uuid.New()
	payload, err := jwtS.CreateToken(userID, time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, payload.Subject, userID.String())
	assert.False(t, payload.ExpiresAt.Before(time.Now()))
	assert.Less(t, 0, len(payload.Token))

	// assert that tokenString was created with the provided key
	_, err = jwt.Parse(payload.Token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
	assert.NoError(t, err)

	// assert that jwtS can verify the token as well
	_, err = jwtS.Verify(payload.Token + string(mustGenerateRandomKey(10)))
	assert.Error(t, err)

	payload2, err := jwtS.Verify(payload.Token)
	assert.NoError(t, err)
	assert.Equal(t, payload, payload2)

	// assert that verify doesn't verify expired tokens
	payload, err = jwtS.CreateToken(userID, -1*time.Minute)
	assert.NoError(t, err, "CreateToken should not return an error")
	_, err = jwtS.Verify(payload.Token)
	assert.Error(t, err, "Verify should return an error for expired token")
}

func mustGenerateRandomKey(keyLength int) []byte {
	key := make([]byte, keyLength)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return key
}

// key, err := generateRandomKey(32)
// assert.NoError(t, err)
// jwtSvc := jwt.NewJWT(key)

// testCases := []struct {
// 	name       string
// 	tokenS     func() string
// 	assertions func(t *testing.T, err error)
// }{
// 	{
// 		name:   "Invalid Token",
// 		tokenS: func() string { return "invalidToken" },
// 		assertions: func(t *testing.T, err error) {
// 			assert.ErrorContains(t, err, "verify: error parsing tokenString")
// 		},
// 	},
// 	{
// 		name: "False secret",
// 		tokenS: func() string {
// 			token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, jwtLib.MapClaims{})
// 			tokenS, err := token.SignedString([]byte("false secret"))
// 			assert.NoError(t, err)
// 			return tokenS
// 		},
// 		assertions: func(t *testing.T, err error) {
// 			assert.ErrorContains(t, err, "verify: error parsing tokenString")
// 		},
// 	},
// 	{
// 		name: "Expired Token",
// 		tokenS: func() string {
// 			oneMinuteAgo := jwtLib.NewNumericDate(time.Now().Add(-1 * time.Minute))
// 			claims := setupClaims(oneMinuteAgo, "", "")
// 			return setupJwtTokenString(t, claims, j.key)
// 		},
// 		assertions: func(t *testing.T, err error) {
// 			assert.ErrorContains(t, err, "verify: ")
// 		},
// 	},
// 	// {
// 	// 	name: "No expiration date set",
// 	// 	tokenS: func() string {

// 	// 		token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, jwtLib.MapClaims{})
// 	// 		tokenS, err := token.SignedString([]byte(j.key))
// 	// 		assert.NoError(t, err)
// 	// 		return tokenS
// 	// 	},
// 	// 	assertions: func(t *testing.T, err error) {
// 	// 		assert.EqualError(t, err, "authenticate: no expiration date set")
// 	// 	},
// 	// },
// 	// {
// 	// 	name: "False issuer",
// 	// 	bearerToken: func() string {
// 	// 		inOneHour := jwtLib.NewNumericDate(time.Now().Add(1 * time.Hour))
// 	// 		claims := setupClaims(inOneHour, "false issuer", "")
// 	// 		return setupJwtTokenString(t, claims, secret)
// 	// 	},
// 	// 	assertion: func(t *testing.T, err error) {
// 	// 		assert.EqualError(t, err, "authenticate: incorrect issuer")
// 	// 	},
// 	// },
// 	//
// 	// {
// 	// 	name:        "Wrong method",
// 	// 	bearerToken: func() string { return getBearerTokenEcdsa256(t) },
// 	// 	assertion: func(t *testing.T, err error) {
// 	// 		assert.ErrorContains(t, err, "authenticate: verify: ")
// 	// 	},
// 	// },
// }

// for _, tc := range testCases {
// 	_, err := j.Verify(tc.tokenS())
// 	tc.assertions(t, err)
// }
//
// func setupJwtTokenString(t *testing.T, claims jwtLib.MapClaims, secret string) string {
// 	t.Helper()
// 	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
// 	tokenS, err := token.SignedString([]byte(secret))
// 	assert.NoError(t, err)
// 	bearerToken := "Bearer " + tokenS
// 	return bearerToken
// }

// func setupClaims(exp *jwtLib.NumericDate, iss, sub string) jwtLib.MapClaims {
// 	claims := jwtLib.MapClaims{
// 		"exp": exp,
// 		"iss": iss,
// 		"sub": sub,
// 	}
// 	return claims
// }
