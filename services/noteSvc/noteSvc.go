// Package noteSvc wraps the data/store layer
// handles Crud operations on note domain type
// make changes persistent by calling data/store layer
package noteSvc

import (
	"fmt"

	"github.com/Keisn1/note-taking-app/domain/note"
	"github.com/google/uuid"
)

type NotesService struct {
	notes noteRepo
}

type noteRepo interface {
	Delete(noteID uuid.UUID) error
	Create(n note.Note) error
	Update(note note.Note) error
	GetNoteByID(noteID uuid.UUID) (note.Note, error)
	GetNotesByUserID(userID uuid.UUID) ([]note.Note, error)
}

func NewNotesService(nR noteRepo) NotesService {
	return NotesService{notes: nR}
}

func (ns NotesService) Delete(noteID uuid.UUID) error {
	err := ns.notes.Delete(noteID)
	if err != nil {
		return fmt.Errorf("delete: [%s]", noteID)
	}
	return nil
}

func (ns NotesService) Create(nN note.NewNote) (note.Note, error) {
	n := note.MakeNoteFromNewNote(nN)
	err := ns.notes.Create(n)
	if err != nil {
		return note.Note{}, err
	}
	return n, nil
}

func (ns NotesService) Update(n, newN note.Note) (note.Note, error) {
	if !newN.GetTitle().IsEmpty() {
		n.SetTitle(newN.GetTitle().String())
	}

	if !newN.GetContent().IsEmpty() {
		n.SetContent(newN.GetContent().String())
	}

	err := ns.notes.Update(n)
	if err != nil {
		return note.Note{}, fmt.Errorf("update: %w", err)
	}
	return n, nil
}

func (nS NotesService) GetNoteByID(noteID uuid.UUID) (note.Note, error) {
	n, err := nS.notes.GetNoteByID(noteID)
	if err != nil {
		return note.Note{}, fmt.Errorf("getNoteByID: [%s]: %w", noteID, err)
	}
	return n, nil
}

func (nS NotesService) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) {
	notes, err := nS.notes.GetNotesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("getNoteByUserID: [%s]: %w", userID, err)
	}
	return notes, nil
}
