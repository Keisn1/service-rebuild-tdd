// Package usernoteSvc wraps the data/store layer
// handles Crud operations on aggregate usernote
// make changes persistent by calling data/store layer
package usernoteSvc

import (
	"fmt"
	"github.com/google/uuid"
)

type Note struct {
	NoteID  uuid.UUID
	Title   string
	Content string
	UserID  uuid.UUID
}

type notesRepo struct {
	notes map[uuid.UUID]Note
}

func NewNotesRepo(notes []Note) (notesRepo, error) {
	var nR notesRepo
	if err := noDuplicate(notes); err != nil {
		return notesRepo{}, fmt.Errorf("newNotesRepo: %w", err)
	}

	nR.notes = make(map[uuid.UUID]Note)
	for _, n := range notes {
		nR.notes[n.NoteID] = n
	}
	return nR, nil
}

func (nR notesRepo) Update(noteID uuid.UUID, newTitle string) Note {
	n := nR.notes[noteID]
	n.Title = newTitle
	nR.notes[noteID] = n
	return n
}

func (nR notesRepo) GetNoteByID(noteID uuid.UUID) Note {
	for _, n := range nR.notes {
		if n.NoteID == noteID {
			return n
		}
	}
	return Note{}
}

func (nR notesRepo) GetNotesByUserID(userID uuid.UUID) []Note {
	var ret []Note
	for _, n := range nR.notes {
		if n.UserID == userID {
			n.NoteID = uuid.UUID{0}
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
