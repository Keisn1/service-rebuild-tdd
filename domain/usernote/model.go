package usernote

import (
	"github.com/Keisn1/note-taking-app/domain/entities"
)

type UserNote struct {
	note *entities.Note
	user *entities.Person
}
