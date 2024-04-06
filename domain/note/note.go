package note

import (
	"github.com/google/uuid"
)

type Note struct {
	NoteID  uuid.UUID
	Title   Title
	Content Content
	UserID  uuid.UUID
}

type Title *string
