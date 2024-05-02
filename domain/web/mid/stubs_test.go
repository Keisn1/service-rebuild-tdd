package mid_test

import (
	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/google/uuid"
)

type StubNoteService struct {
	notes map[uuid.UUID]note.Note
}

func (ns StubNoteService) Delete(noteID uuid.UUID) error                { return nil }
func (ns StubNoteService) Create(nN note.UpdateNote) (note.Note, error) { return note.Note{}, nil }
func (ns StubNoteService) Update(n note.Note, newN note.UpdateNote) (note.Note, error) {
	return note.Note{}, nil
}
func (ns StubNoteService) GetNoteByID(noteID uuid.UUID) (note.Note, error) {
	return ns.notes[noteID], nil
}
func (ns StubNoteService) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) { return nil, nil }
