// Package usernoteSvc wraps the data/store layer
// handles Crud operations on aggregate usernote
// make changes persistent by calling data/store layer
package usernoteSvc

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type Note struct {
	NoteID  uuid.UUID
	Title   string
	Content string
	UserID  uuid.UUID
}

type NotesRepo struct {
	notes map[uuid.UUID]Note
}

func NewNotesRepo(notes []Note) (NotesRepo, error) {
	var nR NotesRepo
	if err := noDuplicate(notes); err != nil {
		return NotesRepo{}, fmt.Errorf("newNotesRepo: %w", err)
	}

	nR.notes = make(map[uuid.UUID]Note)
	for _, n := range notes {
		nR.notes[n.NoteID] = n
	}
	return nR, nil
}

type NoteService struct {
	notes NotesRepo
}

func NewNoteService(nRepo NotesRepo) NoteService {
	return NoteService{
		notes: nRepo,
	}
}

func (ns NoteService) GetNoteByID(noteID uuid.UUID, userID uuid.UUID) (Note, error) {
	n, _ := ns.notes.GetNoteByID(noteID, userID)

	if n.NoteID == noteID {
		if n.UserID != userID {
			return Note{}, fmt.Errorf("getNoteByID: user unauthorized")
		}
		return n, nil
	}
	return Note{}, nil
}

func (nR NotesRepo) GetNoteByID(noteID uuid.UUID, userID uuid.UUID) (Note, error) {
	for _, n := range nR.notes {
		if n.NoteID == noteID {
			return n, nil
		}
	}
	return Note{}, errors.New("")
}

func (nR NotesRepo) GetNotesByUserID(userID uuid.UUID) []Note {
	var ret []Note
	for _, n := range nR.notes {
		if n.UserID == userID {
			ret = append(ret, n)
		}
	}
	return ret
}

func noDuplicate(notes []Note) error {
	noteIDSet := make(map[uuid.UUID]struct{})
	for _, n := range notes {
		if _, ok := noteIDSet[n.NoteID]; ok {
			return fmt.Errorf("duplicate noteID [%s]", n.NoteID)
		}
		noteIDSet[n.NoteID] = struct{}{}
	}
	return nil
}
