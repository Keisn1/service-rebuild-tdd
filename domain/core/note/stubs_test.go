package note_test

import (
	"errors"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/google/uuid"
)

type StubNoteRepo struct {
	notes map[uuid.UUID]note.Note
}

func (nR StubNoteRepo) Delete(noteID uuid.UUID) error { return nil }

func (nR StubNoteRepo) Create(n note.Note) error {
	if n.GetTitle().String() == "invalid title" {
		return errors.New("error in noteRepo")
	}
	return nil
}

func (nR StubNoteRepo) Update(note note.Note) error {
	return nil
}

func (nR StubNoteRepo) GetNoteByID(noteID uuid.UUID) (note.Note, error) {
	return note.Note{}, nil
}

func (nR StubNoteRepo) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) {
	return nil, nil
}
