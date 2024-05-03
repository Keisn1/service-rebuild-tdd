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

type UpdateUser struct {
	Name     string
	Email    mail.Address
	Password string
}
