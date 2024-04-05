package user

import (
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByID(userID uuid.UUID) (User, error)
}
