package user_test

import (
	"context"
	"net/mail"
	"testing"

	"github.com/Keisn1/note-taking-app/domain/core/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_QueryByID(t *testing.T) {
	t.Run("I can get a user by the ID", func(t *testing.T) {
		users := []user.User{
			{ID: uuid.UUID{1}, Name: "rob", Email: mail.Address{Address: "rob@example.com"}},
			{ID: uuid.UUID{2}, Name: "anna", Email: mail.Address{Address: "anna@example.com"}},
		}
		userSvc := user.NewUserSvc(user.NewRepo(users))

		type testCase struct {
			name      string
			ctx       context.Context
			userID    uuid.UUID
			want      user.User
			wantError bool
		}

		testCases := []testCase{
			{
				name:   "retrieve rob",
				userID: uuid.UUID{1},
				want: user.User{
					ID:    uuid.UUID{1},
					Name:  "rob",
					Email: mail.Address{Address: "rob@example.com"},
				},
			},
			{
				name:   "retrieve anna",
				userID: uuid.UUID{2},
				want: user.User{
					ID:    uuid.UUID{2},
					Name:  "anna",
					Email: mail.Address{Address: "anna@example.com"},
				},
			},
			{
				name:      "return error on missing user",
				userID:    uuid.New(),
				want:      user.User{},
				wantError: true,
			},
		}
		for _, tc := range testCases {
			got, err := userSvc.QueryByID(tc.ctx, tc.userID)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want, got)
		}
	})
}
