package domain_test

import (
	"errors"
	"fmt"
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

func (sUNR *StubUserNoteRepository) GetNotesByUserID(uID uuid.UUID) ([]usernote.UserNote, error) {
	var ret []usernote.UserNote
	for _, n := range sUNR.usernotes {
		if n.GetUserID() == uID {
			ret = append(ret, n)
		}
	}
	return ret, nil
}

func TestService(t *testing.T) {
	uID1 := uuid.UUID([16]byte{1})
	nID1 := uuid.UUID([16]byte{2})
	nID2 := uuid.UUID([16]byte{3})

	uID2 := uuid.UUID([16]byte{4})
	nID3 := uuid.UUID([16]byte{5})
	nID4 := uuid.UUID([16]byte{6})

	note1 := usernote.NewUserNote(nID1, "", "", uID1)
	note2 := usernote.NewUserNote(nID2, "", "", uID1)
	note3 := usernote.NewUserNote(nID3, "", "", uID2)
	note4 := usernote.NewUserNote(nID4, "", "", uID2)

	u := &StubUserNoteRepository{
		usernotes: map[uuid.UUID]usernote.UserNote{
			nID1: note1, nID2: note2, nID3: note3, nID4: note4,
		},
	}
	s := domain.NewService(domain.WithUserNoteRepository(u))

	t.Run("Get note by noteID", func(t *testing.T) {
		want := note1
		got, err := s.QueryByID(nID1)
		assert.NoError(t, err)
		assert.Equal(t, want, got)

		want = note2
		got, err = s.QueryByID(nID2)
		assert.Equal(t, want, got)
		assert.NoError(t, err)
	})

	t.Run("Return error for missing note", func(t *testing.T) {
		nIDx := uuid.UUID([16]byte{100})
		_, err := s.QueryByID(nIDx)
		expectedErrorSubString := fmt.Sprintf("querybyid: noteID[%s]", nIDx)
		assert.ErrorContains(t, err, expectedErrorSubString)
	})

	t.Run("Return notes by UserID", func(t *testing.T) {
		wantNotes := []usernote.UserNote{note1, note2}
		gotNotes, err := s.QueryByUserID(uID1)
		assert.NoError(t, err)
		assert.Equal(t, wantNotes, gotNotes)

		wantNotes = []usernote.UserNote{note3, note4}
		gotNotes, err = s.QueryByUserID(uID2)
		assert.NoError(t, err)
		assert.Equal(t, wantNotes, gotNotes)
	})

	t.Run("Return empty notes if no notes found for user with UserID", func(t *testing.T) {
		uIDx := uuid.UUID([16]byte{100})
		wantNotes := []usernote.UserNote{}
		gotNotes, err := s.QueryByUserID(uIDx)
		assert.NoError(t, err)
		assert.Equal(t, wantNotes, gotNotes)
	})
}
