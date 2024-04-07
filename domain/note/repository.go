package note

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type NoteRepo struct {
	notes map[uuid.UUID]Note
}

func NewNotesRepo(notes []Note) (NoteRepo, error) {
	var nR NoteRepo
	if err := noDuplicate(notes); err != nil {
		return NoteRepo{}, fmt.Errorf("newNotesRepo: %w", err)
	}

	nR.notes = make(map[uuid.UUID]Note)
	for _, n := range notes {
		nR.notes[n.GetID()] = n
	}
	return nR, nil
}

func (nR NoteRepo) Create(n Note) {
	nR.notes[n.GetID()] = n
}

func (nR NoteRepo) Update(note Note) error {
	if _, ok := nR.notes[note.GetID()]; ok {
		nR.notes[note.GetID()] = note
		return nil
	}
	return errors.New("")
}

func (nR NoteRepo) GetNoteByID(noteID uuid.UUID) (Note, error) {
	for _, n := range nR.notes {
		if n.GetID() == noteID {
			return n, nil
		}
	}
	return Note{}, fmt.Errorf("GetNoteByID: Not found [%s]", noteID)
}

func (nR NoteRepo) GetNotesByUserID(userID uuid.UUID) ([]Note, error) {
	var ret []Note
	var found bool
	for _, n := range nR.notes {
		if n.GetUserID() == userID {
			found = true
			n.SetID(uuid.UUID{0})
			ret = append(ret, n)
		}
	}
	if !found {
		return nil, fmt.Errorf("getNotesByUserID: not found [%s]", userID)
	}
	return ret, nil
}

func noDuplicate(notes []Note) error {
	noteIDSet := make(map[uuid.UUID]struct{})
	for _, n := range notes {
		if _, ok := noteIDSet[n.GetID()]; ok {
			return fmt.Errorf("duplicate noteID [%s]", n.GetID())
		}
		noteIDSet[n.GetID()] = struct{}{}
	}
	return nil
}
