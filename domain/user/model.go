package user

import (
	"net/mail"

	"github.com/Keisn1/note-taking-app/domain/entities"
	"github.com/google/uuid"
)

type User struct {
	person *entities.Person
}

func NewUser(name, email string) User {
	return User{
		person: &entities.Person{
			ID:    uuid.New(),
			Name:  entities.Username(name),
			Email: mail.Address{Address: email},
		},
	}
}

func (u User) GetID() uuid.UUID {
	return u.person.ID
}
