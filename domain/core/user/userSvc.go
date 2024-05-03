package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
}

type Repo struct {
	users map[uuid.UUID]User
}

func NewRepo(users []User) Repo {
	us := make(map[uuid.UUID]User)
	for _, u := range users {
		us[u.ID] = u
	}
	return Repo{users: us}
}

type Svc struct {
	repo Repo
}

func NewUserSvc(repo Repo) Service {
	return Svc{repo: repo}
}

func (us Svc) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	if user, ok := us.repo.users[userID]; ok {
		return user, nil
	}
	return User{}, errors.New("user not found")
}
