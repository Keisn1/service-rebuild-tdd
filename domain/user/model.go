package user

import (
	"github.com/Keisn1/note-taking-app/domain/entities"
	"github.com/google/uuid"
)

type User struct {
	person *entities.Person
	notes  []*entities.Note
}

func NewUser(name, email string) User {
	return User{
		person: &entities.Person{
			ID:    uuid.New(),
			Name:  entities.Username(name),
			Email: entities.Email(email),
		},
		notes: []*entities.Note{},
	}
}

func (u User) GetID() uuid.UUID {
	return u.person.ID

}
