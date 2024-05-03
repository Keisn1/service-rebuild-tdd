package memory

import (
	"context"
	"errors"
	"fmt"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/google/uuid"
)

type Repo struct {
	notes map[uuid.UUID]note.Note
}

func NewRepo(notes []note.Note) (Repo, error) {
	var nR Repo
	if err := noDuplicate(notes); err != nil {
		return Repo{}, fmt.Errorf("newNotesRepo: %w", err)
	}

	nR.notes = make(map[uuid.UUID]note.Note)
	for _, n := range notes {
		nR.notes[n.GetID()] = n
	}
	return nR, nil
}

func MustNewRepo(notes []note.Note) Repo {
	nr, err := NewRepo(notes)
	if err != nil {
		panic(err)
	}
	return nr
}

func (nR Repo) Delete(noteID uuid.UUID) error {
	if _, ok := nR.notes[noteID]; ok {
		delete(nR.notes, noteID)
		return nil
	}
	return fmt.Errorf("delete: not found [%s]", noteID)

}

func (nR Repo) Create(n note.Note) error {
	if _, ok := nR.notes[n.GetID()]; ok {
		return fmt.Errorf("create: already present %s", n.GetID())
	}
	nR.notes[n.GetID()] = n
	return nil
}

func (nR Repo) Update(note note.Note) error {
	if _, ok := nR.notes[note.GetID()]; ok {
		nR.notes[note.GetID()] = note
		return nil
	}
	return errors.New("")
}

func (nR Repo) QueryByID(ctx context.Context, noteID uuid.UUID) (note.Note, error) {
	for _, n := range nR.notes {
		if n.GetID() == noteID {
			return n, nil
		}
	}
	return note.Note{}, fmt.Errorf("GetNoteByID: Not found [%s]", noteID)
}

func (nR Repo) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) {
	var ret []note.Note
	var found bool
	for _, n := range nR.notes {
		if n.GetUserID() == userID {
			found = true
			ret = append(ret, n)
		}
	}
	if !found {
		return nil, fmt.Errorf("getNotesByUserID: not found [%s]", userID)
	}
	return ret, nil
}

func noDuplicate(notes []note.Note) error {
	noteIDSet := make(map[uuid.UUID]struct{})
	for _, n := range notes {
		if _, ok := noteIDSet[n.GetID()]; ok {
			return fmt.Errorf("duplicate noteID [%s]", n.GetID())
		}
		noteIDSet[n.GetID()] = struct{}{}
	}
	return nil
}
