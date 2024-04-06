// Package usernoteSvc wraps the data/store layer
// handles Crud operations on aggregate usernote
// make changes persistent by calling data/store layer
package usernoteSvc

import (
	"fmt"

	"github.com/Keisn1/note-taking-app/domain/note"
	"github.com/google/uuid"
)

type NotesService struct {
	notes note.NoteRepo
}

func NewNotesService(nR note.NoteRepo) NotesService {
	return NotesService{notes: nR}
}

func (ns NotesService) Update(note, newNote note.Note) error {
	// TODO: anything calling into the service, shall already talk the language of the service => Notes

	if !newNote.GetTitle().IsEmpty() {
		note.GetTitle().Set(newNote.GetTitle().Get())
	}

	if !newNote.GetContent().IsEmpty() {
		note.GetContent().Set(newNote.GetContent().Get())
	}

	err := ns.notes.Update(note)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (nS NotesService) GetNoteByID(noteID uuid.UUID) note.Note {
	return nS.notes.GetNoteByID(noteID)
}

func (nS NotesService) GetNotesByUserID(userID uuid.UUID) []note.Note {
	return nS.notes.GetNotesByUserID(userID)
}
