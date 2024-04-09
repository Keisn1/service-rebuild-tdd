package note

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNoteNotFound = errors.New("the product was not found")
)

type NoteRepo interface {
	Delete(noteID uuid.UUID) error
	Create(n Note) error
	Update(note Note) error
	GetNoteByID(noteID uuid.UUID) (Note, error)
	GetNotesByUserID(userID uuid.UUID) ([]Note, error)
}
