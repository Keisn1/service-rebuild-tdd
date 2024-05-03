package note

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNoteNotFound = errors.New("the note was not found")
)

type Repo interface {
	Delete(noteID uuid.UUID) error
	Create(n Note) error
	Update(note Note) error
	QueryByID(ctx context.Context, noteID uuid.UUID) (Note, error)
	GetNotesByUserID(userID uuid.UUID) ([]Note, error)
}
