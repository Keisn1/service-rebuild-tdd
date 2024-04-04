package user

import "github.com/Keisn1/note-taking-app/domain/entities"

type User struct {
	person *entities.Person
	notes  []*entities.Note
}
