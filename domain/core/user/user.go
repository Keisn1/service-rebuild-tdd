package user

import (
	"net/mail"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         Name
	Email        Email
	PasswordHash []byte
}

type UpdateUser struct {
	Name     Name
	Email    Email
	Password string
}

type Name struct {
	name *string
}

func NewName(un string) Name {
	return Name{name: &un}
}

func (u *User) GetName() Name       { return u.Name }
func (u *User) GetID() uuid.UUID    { return u.ID }
func (u *UpdateUser) GetName() Name { return u.Name }

func (n Name) IsEmpty() bool   { return n.name == nil }
func (n Name) Set(name string) { *n.name = name }
func (n Name) String() string {
	if n.IsEmpty() {
		return ""
	}
	return *n.name
}

type Email struct {
	email *mail.Address
}

func NewEmail(email mail.Address) Email {
	return Email{email: &email}
}
func (e Email) IsEmpty() bool          { return e.email == nil }
func (e Email) Set(email mail.Address) { *e.email = email }
func (e Email) String() mail.Address {
	if e.IsEmpty() {
		return mail.Address{}
	}
	return *e.email
}
