package note

import (
	"fmt"

	"github.com/google/uuid"
)

type NotesServiceI interface {
	Delete(noteID uuid.UUID) error
	Create(nN UpdateNote) (Note, error)
	Update(n Note, newN UpdateNote) (Note, error)
	GetNoteByID(noteID uuid.UUID) (Note, error)
	GetNotesByUserID(userID uuid.UUID) ([]Note, error)
}

type NotesService struct {
	notes NoteRepo
}

func NewNotesService(nR NoteRepo) NotesService {
	return NotesService{notes: nR}
}

func (ns NotesService) Delete(noteID uuid.UUID) error {
	err := ns.notes.Delete(noteID)
	if err != nil {
		return fmt.Errorf("delete: [%s]", noteID)
	}
	return nil
}

func (ns NotesService) Create(nN UpdateNote) (Note, error) {
	n := Note{
		NoteID:  uuid.New(),
		Title:   nN.Title,
		Content: nN.Content,
		UserID:  nN.UserID,
	}

	err := ns.notes.Create(n)
	if err != nil {
		return Note{}, err
	}
	return n, nil
}

func (ns NotesService) Update(n Note, newN UpdateNote) (Note, error) {
	if !newN.GetTitle().IsEmpty() {
		n.SetTitle(newN.GetTitle().String())
	}

	if !newN.GetContent().IsEmpty() {
		n.SetContent(newN.GetContent().String())
	}

	err := ns.notes.Update(n)
	if err != nil {
		return Note{}, fmt.Errorf("update: %w", err)
	}
	return n, nil
}

func (nS NotesService) GetNoteByID(noteID uuid.UUID) (Note, error) {
	n, err := nS.notes.GetNoteByID(noteID)
	if err != nil {
		return Note{}, fmt.Errorf("getNoteByID: [%s]: %w", noteID, err)
	}
	return n, nil
}

func (nS NotesService) GetNotesByUserID(userID uuid.UUID) ([]Note, error) {
	notes, err := nS.notes.GetNotesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("getNoteByUserID: [%s]: %w", userID, err)
	}
	return notes, nil
}
