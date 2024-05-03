package memory

import (
	"context"
	"errors"

	"github.com/Keisn1/note-taking-app/domain/core/user"
	"github.com/google/uuid"
)

type InMemoryRepo struct {
	users map[uuid.UUID]user.User
}

func NewRepo(users []user.User) InMemoryRepo {
	us := make(map[uuid.UUID]user.User)
	for _, u := range users {
		us[u.ID] = u
	}
	return InMemoryRepo{users: us}
}

func (r InMemoryRepo) Update(ctx context.Context, u user.User) error {
	r.users[u.ID] = u
	return nil
}

func (r InMemoryRepo) Create(ctx context.Context, u user.User) error {
	r.users[u.ID] = u
	return nil
}

func (r InMemoryRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	if _, ok := r.users[userID]; !ok {
		return errors.New("user not found")
	}
	delete(r.users, userID)
	return nil
}

func (r InMemoryRepo) QueryByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	if user, ok := r.users[userID]; ok {
		return user, nil
	}
	return user.User{}, errors.New("user not found")
}
