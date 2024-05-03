package user

import (
	"context"

	"github.com/google/uuid"
)

type Repo interface {
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	Create(ctx context.Context, u User) error
	Update(ctx context.Context, u User) error
	Delete(ctx context.Context, userID uuid.UUID) error
}
