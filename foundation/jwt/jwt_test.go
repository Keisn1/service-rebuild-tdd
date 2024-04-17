package jwt_test

import (
	"os"
	"testing"

	"github.com/Keisn1/note-taking-app/foundation/jwt"
	jwtLib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWT_CreateToken(t *testing.T) {
	j := jwt.NewJWT(os.Getenv("JWT_SECRET_KEY"))

	userID := uuid.New()
	jwtPayload := j.CreateToken(userID)
	assert.Equal(t, jwt.Payload.UserID, userID)
}

func TestJWT_Verify(t *testing.T) {
	j := jwt.NewJWT(os.Getenv("JWT_SECRET_KEY"))

	testCases := []struct {
		name       string
		tokenS     func() string
		assertions func(t *testing.T, err error)
	}{
		{
			name:   "Invalid Token",
			tokenS: func() string { return "invalidToken" },
			assertions: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "verify: error parsing tokenString")
			},
		},
		{
			name: "False secret",
			tokenS: func() string {
				token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, jwtLib.MapClaims{})
				tokenS, err := token.SignedString([]byte("false secret"))
				assert.NoError(t, err)
				return tokenS
			},
			assertions: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "verify: error parsing tokenString")
			},
		},
		// {
		// 	name: "Expired Token",
		// 	tokenS: func() string {
		// 		oneMinuteAgo := jwtLib.NewNumericDate(time.Now().Add(-1 * time.Minute))
		// 		claims := setupClaims(oneMinuteAgo, "", "")

		// 		return setupJwtTokenString(t, claims, j.key)
		// 	},
		// 	assertions: func(t *testing.T, err error) {
		// 		assert.ErrorContains(t, err, "authenticate: error parsing tokenString")
		// 	},
		// },
		// {
		// 	name: "No expiration date set",
		// 	tokenS: func() string {

		// 		token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, jwtLib.MapClaims{})
		// 		tokenS, err := token.SignedString([]byte(j.key))
		// 		assert.NoError(t, err)
		// 		return tokenS
		// 	},
		// 	assertions: func(t *testing.T, err error) {
		// 		assert.EqualError(t, err, "authenticate: no expiration date set")
		// 	},
		// },
	}

	for _, tc := range testCases {
		_, err := j.Verify(tc.tokenS())
		tc.assertions(t, err)
	}
}
