package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type Service interface {
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	Create(ctx context.Context, nu UpdateUser) (User, error)
	Update(ctx context.Context, u User, uu UpdateUser) (User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}

type Svc struct {
	repo Repo
}

func NewSvc(repo Repo) Service {
	return Svc{repo: repo}
}

func (s Svc) Update(ctx context.Context, u User, newU UpdateUser) (User, error) {
	_, err := s.repo.QueryByID(ctx, u.ID)
	if err != nil {
		return User{}, err
	}

	if !newU.Name.IsEmpty() {
		u.Name = newU.Name
	}

	if !newU.Email.IsEmpty() {
		u.Email = newU.Email
	}

	if !newU.Password.IsEmpty() {
		pwHash, err := bcrypt.GenerateFromPassword([]byte(newU.Password.String()), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("create: %w: %w", ErrInvalidPassword, err)
		}
		u.PasswordHash = pwHash
	}

	s.repo.Update(ctx, u)

	return u, nil
}

func (s Svc) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (s Svc) Create(ctx context.Context, newU UpdateUser) (User, error) {
	if len(newU.Password.String()) == 0 {
		return User{}, fmt.Errorf("create: %w", ErrInvalidPassword)
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(newU.Password.String()), bcrypt.DefaultCost)
	if err != nil {

		return User{}, fmt.Errorf("create: %w: %w", ErrInvalidPassword, err)
	}

	u := User{
		ID:           uuid.New(),
		Name:         newU.Name,
		Email:        newU.Email,
		PasswordHash: pwHash,
	}

	s.repo.Create(ctx, u)
	return u, nil
}

func (s Svc) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	u, err := s.repo.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("queryByID: %w", err)
	}
	return u, nil
}
