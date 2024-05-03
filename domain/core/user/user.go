package user

import (
	"net/mail"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         Name
	Email        mail.Address
	PasswordHash []byte
}

type UpdateUser struct {
	Name     Name
	Email    *mail.Address
	Password string
}

type Name struct {
	userName *string
}

func NewName(un string) Name {
	return Name{userName: &un}
}

func (u *User) GetName() Name { return u.Name }
