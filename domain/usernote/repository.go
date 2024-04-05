package usernote

import (
	"github.com/google/uuid"
)

type UserNoteRepository interface {
	GetNoteByID(noteID uuid.UUID) (UserNote, error)
	GetNotesByUserID(userID uuid.UUID) ([]UserNote, error)
	Create(userID uuid.UUID, title, content string) (UserNote, error)
}
