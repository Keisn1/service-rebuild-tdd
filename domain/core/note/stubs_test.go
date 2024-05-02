package note_test

import (
	"errors"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/google/uuid"
)

type ErrorNoteRepo struct {
	notes map[uuid.UUID]note.Note
}

func (nR ErrorNoteRepo) Create(n note.Note) error                               { return errors.New("error in noteRepo") }
func (nR ErrorNoteRepo) Delete(noteID uuid.UUID) error                          { return nil }
func (nR ErrorNoteRepo) Update(note note.Note) error                            { return nil }
func (nR ErrorNoteRepo) GetNoteByID(noteID uuid.UUID) (note.Note, error)        { return note.Note{}, nil }
func (nR ErrorNoteRepo) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) { return nil, nil }
