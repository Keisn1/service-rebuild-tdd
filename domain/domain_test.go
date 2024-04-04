package domain_test

import (
	"errors"
	"testing"

	"github.com/Keisn1/note-taking-app/domain"
	"github.com/Keisn1/note-taking-app/domain/usernote"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type StubUserNoteRepository struct {
	usernotes map[uuid.UUID]usernote.UserNote
}

func (sUNR *StubUserNoteRepository) GetNoteByID(nID uuid.UUID) (usernote.UserNote, error) {
	n, ok := sUNR.usernotes[nID]
	if !ok {
		return usernote.UserNote{}, errors.New("error")
	}
	return n, nil
}

func TestNotes(t *testing.T) {
	nID1 := uuid.UUID([16]byte{1})
	nID2 := uuid.UUID([16]byte{2})
	uID1 := uuid.UUID([16]byte{3})
	uID2 := uuid.UUID([16]byte{4})
	note1 := usernote.NewUserNote(nID1, "", "", uID1)
	note2 := usernote.NewUserNote(nID2, "", "", uID2)
	u := &StubUserNoteRepository{
		usernotes: map[uuid.UUID]usernote.UserNote{
			nID1: note1,
			nID2: note2,
		},
	}
	s := domain.NewService(domain.WithUserNoteRepository(u))

	t.Run("Return note for noteID", func(t *testing.T) {
		want := note1
		got, err := s.GetNoteByID(nID1)
		assert.Equal(t, want, got)
		assert.NoError(t, err)

		want = note2
		got, err = s.GetNoteByID(nID2)
		assert.Equal(t, want, got)
		assert.NoError(t, err)
	})

	t.Run("Return Error for missing note", func(t *testing.T) {
		nIDx := uuid.UUID([16]byte{100})
		_, err := s.GetNoteByID(nIDx)
		assert.Error(t, err)
	})
}
