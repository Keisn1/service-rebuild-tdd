package auth_test

import (
	"testing"

	"github.com/Keisn1/note-taking-app/domain/web/auth"
	"github.com/Keisn1/note-taking-app/foundation/common"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	key := common.MustGenerateRandomKey(32)
	a := auth.NewAuth(auth.MustNewJWTService(key))

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
			bearerToken: func() string { return "Bearer asdf;ljasdfl;j" },
			assertion: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "authenticate: verify: ")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := a.Authenticate(tc.bearerToken())
			tc.assertion(t, err)
		})
	}
}
