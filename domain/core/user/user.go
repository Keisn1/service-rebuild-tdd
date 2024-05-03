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
	Password Password
}

type Name struct {
	name *string
}

func NewName(un string) Name {
	return Name{name: &un}
}

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

func NewEmail(email string) Email {
	return Email{email: &mail.Address{Address: email}}
}
func (e Email) IsEmpty() bool          { return e.email == nil }
func (e Email) Set(email mail.Address) { *e.email = email }
func (e Email) String() mail.Address {
	if e.IsEmpty() {
		return mail.Address{}
	}
	return *e.email
}

type Password struct {
	password *string
}

func NewPassword(un string) Password {
	return Password{password: &un}
}

func (p Password) IsEmpty() bool       { return p.password == nil }
func (p Password) Set(password string) { *p.password = password }
func (p Password) String() string {
	if p.IsEmpty() {
		return ""
	}
	return *p.password
}
