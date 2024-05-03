package user_test

import (
	"context"
	"net/mail"
	"testing"

	"github.com/Keisn1/note-taking-app/domain/core/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func Test_Update(t *testing.T) {
	t.Run("Update Name", func(t *testing.T) {
		u := user.User{ID: uuid.UUID{1}, Name: user.NewName("rob"), Email: user.NewEmail(mail.Address{Address: "rob@example.com"})}

		users := []user.User{u}
		svc := user.NewSvc(user.NewRepo(users))

		uu := user.UpdateUser{
			Name: user.NewName("robbie"),
		}

		updatedU, err := svc.Update(context.Background(), u, uu)
		assert.NoError(t, err)
		assert.Equal(t, updatedU.Name, uu.Name)
		assert.Equal(t, updatedU.Name, uu.Name)

		assert.Equal(t, updatedU.Email, u.Email)
		assert.Equal(t, updatedU.Email, u.Email)
	})

	t.Run("Update Email", func(t *testing.T) {
		u := user.User{ID: uuid.UUID{1}, Name: user.NewName("rob"), Email: user.NewEmail(mail.Address{Address: "rob@example.com"})}

		users := []user.User{u}
		svc := user.NewSvc(user.NewRepo(users))

		uu := user.UpdateUser{Email: user.NewEmail(mail.Address{Address: "robbie@example.com"})}

		updatedU, err := svc.Update(context.Background(), u, uu)
		assert.NoError(t, err)
		assert.Equal(t, updatedU.Email, uu.Email)
		assert.Equal(t, updatedU.Email, uu.Email)

		assert.Equal(t, updatedU.Name, u.Name)
		assert.Equal(t, updatedU.Name, u.Name)
	})
}

func Test_Create(t *testing.T) {
	svc := user.NewSvc(user.NewRepo([]user.User{}))

	t.Run("Happy paths", func(t *testing.T) {
		type testCase struct {
			newUser  user.UpdateUser
			wantUser user.User
		}

		testCases := []testCase{
			{
				newUser: user.UpdateUser{
					Name:     user.NewName("rob"),
					Email:    user.NewEmail(mail.Address{Address: "rob@example.com"}),
					Password: "password",
				},
				wantUser: user.User{
					Name:  user.NewName("rob"),
					Email: user.NewEmail(mail.Address{Address: "rob@example.com"}),
				},
			},
			{
				newUser: user.UpdateUser{
					Name:     user.NewName("anna"),
					Email:    user.NewEmail(mail.Address{Address: "anna@example.com"}),
					Password: "passwordAnna",
				},
				wantUser: user.User{
					Name:  user.NewName("anna"),
					Email: user.NewEmail(mail.Address{Address: "anna@example.com"}),
				},
			},
		}

		for _, tc := range testCases {
			createdUser, err := svc.Create(context.Background(), tc.newUser)
			assert.NoError(t, err)
			assert.NotNil(t, createdUser.ID)
			assert.NotEqual(t, createdUser.ID, uuid.UUID{})
			assert.Equal(t, tc.wantUser.Name, createdUser.Name)
			assert.Equal(t, tc.wantUser.Email, createdUser.Email)
			assert.NoError(t, bcrypt.CompareHashAndPassword(createdUser.PasswordHash, []byte(tc.newUser.Password)))

			retrievedUser, err := svc.QueryByID(context.Background(), createdUser.ID)
			assert.NoError(t, err)
			assert.Equal(t, createdUser, retrievedUser)
		}
	})

	t.Run("Password checking", func(t *testing.T) {
		newUser := user.UpdateUser{
			Name:     user.NewName("rob"),
			Email:    user.NewEmail(mail.Address{Address: "rob@example.com"}),
			Password: "",
		}

		_, err := svc.Create(context.Background(), newUser)
		assert.ErrorIs(t, err, user.ErrInvalidPassword)
		assert.ErrorContains(t, err, "create")

		newUser = user.UpdateUser{
			Name:     user.NewName("rob"),
			Email:    user.NewEmail(mail.Address{Address: "rob@example.com"}),
			Password: "72727272727272727272727272727272727272727272727272727272727272727272727272",
		}

		_, err = svc.Create(context.Background(), newUser)
		assert.ErrorIs(t, err, user.ErrInvalidPassword)
		assert.ErrorContains(t, err, "create")
	})
}

func Test_QueryByID(t *testing.T) {
	t.Run("I can get a user by the ID", func(t *testing.T) {
		users := []user.User{
			{ID: uuid.UUID{1}, Name: user.NewName("rob"), Email: user.NewEmail(mail.Address{Address: "rob@example.com"})},
			{ID: uuid.UUID{2}, Name: user.NewName("anna"), Email: user.NewEmail(mail.Address{Address: "anna@example.com"})},
		}
		svc := user.NewSvc(user.NewRepo(users))

		type testCase struct {
			name          string
			userID        uuid.UUID
			want          user.User
			wantError     bool
			errorContains string
		}

		testCases := []testCase{
			{
				name:   "retrieve rob",
				userID: uuid.UUID{1},
				want: user.User{
					ID:    uuid.UUID{1},
					Name:  user.NewName("rob"),
					Email: user.NewEmail(mail.Address{Address: "rob@example.com"}),
				},
			},
			{
				name:   "retrieve anna",
				userID: uuid.UUID{2},
				want: user.User{
					ID:    uuid.UUID{2},
					Name:  user.NewName("anna"),
					Email: user.NewEmail(mail.Address{Address: "anna@example.com"}),
				},
			},
			{
				name:          "return error on missing user",
				userID:        uuid.New(),
				want:          user.User{},
				wantError:     true,
				errorContains: "queryByID",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := svc.QueryByID(context.Background(), tc.userID)
				if tc.wantError {
					assert.Error(t, err)
					assert.ErrorContains(t, err, tc.errorContains)
				} else {
					assert.NoError(t, err)
				}

				assert.Equal(t, tc.want, got)

			})
		}
	})
}
