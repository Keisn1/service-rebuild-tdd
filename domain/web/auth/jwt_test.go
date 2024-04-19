package auth_test

import (
	"testing"
	"time"

	"github.com/Keisn1/note-taking-app/domain/web/auth"
	"github.com/Keisn1/note-taking-app/foundation/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	key := common.MustGenerateRandomKey(24)
	_, err := auth.NewJWTService(key)
	assert.EqualError(t, err, "key minLength 32")

	key = common.MustGenerateRandomKey(32)
	jwtS, err := auth.NewJWTService(key)
	assert.NoError(t, err)

	userID := uuid.New()
	tokenS, err := jwtS.CreateToken(userID, time.Minute)
	assert.NoError(t, err)
	assert.Less(t, 0, len(tokenS))

	// assert that jwtS can verify the token
	claims, err := jwtS.Verify(tokenS)
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), claims.Subject)
	assert.False(t, claims.ExpiresAt.Before(time.Now()))

	// assert that jwtS rejects false the token
	_, err = jwtS.Verify(tokenS + string(common.MustGenerateRandomKey(10)))
	assert.Error(t, err)
	assert.ErrorContains(t, err, "verify: ")

	// assert that verify doesn't verify expired tokens
	tokenS, err = jwtS.CreateToken(userID, -1*time.Minute)
	assert.NoError(t, err, "CreateToken should not return an error")

	_, err = jwtS.Verify(tokenS)
	assert.Error(t, err, "Verify should return an error for expired token")
}
