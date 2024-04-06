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

func (nR NoteRepo) Update(noteID uuid.UUID, newNote Note) error {
	if _, ok := nR.notes[noteID]; ok {
		nR.notes[noteID] = newNote
		return nil
	}
	return errors.New("")
}

func (nR NoteRepo) GetNoteByID(noteID uuid.UUID) Note {
	for _, n := range nR.notes {
		if n.GetID() == noteID {
			return n
		}
	}
	return Note{}
}

func (nR NoteRepo) GetNotesByUserID(userID uuid.UUID) []Note {
	var ret []Note
	for _, n := range nR.notes {
		if n.GetUserID() == userID {
			n.SetID(uuid.UUID{0})
			ret = append(ret, n)
		}
	}
	return ret
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
