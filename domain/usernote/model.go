package usernote

import (
	"github.com/Keisn1/note-taking-app/domain/entities"
	"github.com/google/uuid"
)

type UserNote struct {
	note *entities.Note
	user *entities.Person
}

func NewUserNote(nID uuid.UUID, title, content string, uID uuid.UUID) UserNote {
	return UserNote{
		note: entities.NewNote(nID, title, content),
		user: &entities.Person{ID: uID},
	}
}
