package user

import (
	"net/mail"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID
	Name  string
	Email mail.Address
}
