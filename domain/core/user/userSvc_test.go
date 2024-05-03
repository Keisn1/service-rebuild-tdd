package user_test

import (
	"os/user"
	"testing"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/google/uuid"
)

func Test_QueryByID(t *testing.T) {
	t.Run("I can get a user by the ID", func(t *testing.T) {
		notesS := Setup(t, fixtureUsers())

		type testCase struct {
			userID uuid.UUID
			want   user.User
		}

		testCases := []testCase{
			{
				userID: uuid.UUID{1},
				want: note.User{
					UserID: uuid.UUID{1},
					Name:   "rob",
					Email:  "rob@email.com",
				},
			},
		}
	})
}
