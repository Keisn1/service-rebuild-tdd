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

func (ns NotesService) Create(n note.Note) (note.Note, error) {
	n.SetID(uuid.New())
	ns.notes.Create(n)
	return n, nil
}

func (ns NotesService) Update(n, newN note.Note) (note.Note, error) {
	if !newN.GetTitle().IsEmpty() {
		n.SetTitle(newN.GetTitle().Get())
	}

	if !newN.GetContent().IsEmpty() {
		n.SetContent(newN.GetContent().Get())
	}

	err := ns.notes.Update(n)
	if err != nil {
		return note.Note{}, fmt.Errorf("update: %w", err)
	}
	return n, nil
}

func (nS NotesService) GetNoteByID(noteID uuid.UUID) (note.Note, error) {
	return nS.notes.GetNoteByID(noteID), nil
}

func (nS NotesService) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) {
	return nS.notes.GetNotesByUserID(userID), nil
}
