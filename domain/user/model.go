package user

import (
	"net/mail"

	ents "github.com/Keisn1/note-taking-app/domain/entities"
	"github.com/google/uuid"
)

type User struct {
	person *ents.Person
}

func NewUser(name, email string) User {
	return User{
		person: &ents.Person{
			ID:    ents.UserID(uuid.New()),
			Name:  ents.Username(name),
			Email: mail.Address{Address: email},
		},
	}
}

func (u User) GetID() ents.UserID {
	return u.person.ID
}
