// Package usernoteSvc wraps the data/store layer
// handles Crud operations on aggregate usernote
// make changes persistent by calling data/store layer
package usernoteSvc

import (
	"errors"

	"github.com/Keisn1/note-taking-app/domain/note"
	"github.com/google/uuid"
)

type NotesService struct {
	notes note.NoteRepo
}

func NewNotesService(nR note.NoteRepo) NotesService {
	return NotesService{notes: nR}
}

func (ns NotesService) Update(noteID uuid.UUID, title string) (note.Note, error) {
	// TODO: anything calling into the service, shall already talk the language of the service => Notes
	n := ns.notes.GetNoteByID(noteID)
	x := note.Note{}
	if n == x {
		return note.Note{}, errors.New("update: ")
	}

	ns.notes.Update(noteID, title)
	n.Title = title
	return n, nil
}

func (nS NotesService) GetNoteByID(noteID uuid.UUID) note.Note {
	return nS.notes.GetNoteByID(noteID)
}

func (nS NotesService) GetNotesByUserID(userID uuid.UUID) []note.Note {
	return nS.notes.GetNotesByUserID(userID)
}
