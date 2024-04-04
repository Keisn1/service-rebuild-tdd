package domain_test

import (
	"fmt"
	"testing"

	"github.com/Keisn1/note-taking-app/domain"
	"github.com/Keisn1/note-taking-app/domain/entities"
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
		return usernote.UserNote{}, fmt.Errorf("Note note found")
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
	if len(ret) == 0 {
		return nil, fmt.Errorf("No note found for userID[%s]", uID)
	}
	return ret, nil
}

func TestService(t *testing.T) {
	uID1 := uuid.UUID([16]byte{1})
	uID2 := uuid.UUID([16]byte{2})

	note1 := usernote.NewUserNote("", "", uID1)
	note2 := usernote.NewUserNote("", "", uID1)
	note3 := usernote.NewUserNote("", "", uID2)
	note4 := usernote.NewUserNote("", "", uID2)

	nID1 := note1.GetID()
	nID2 := note2.GetID()
	nID3 := note3.GetID()
	nID4 := note4.GetID()

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

	t.Run("Return error if no notes found for userID", func(t *testing.T) {
		uIDx := uuid.UUID([16]byte{100})
		_, err := s.QueryByUserID(uIDx)
		expectedErrorSubString := fmt.Sprintf("querybyuserid: userID[%s]", uIDx)
		assert.ErrorContains(t, err, expectedErrorSubString)
	})

	t.Run("Add a note", func(t *testing.T) {
		uID := uuid.UUID([16]byte{1})
		got, err := s.Create(uID, "title", "content")
		assert.NoError(t, err)
		assert.Equal(t, got.GetUserID(), uID)
		assert.Equal(t, entities.Title("title"), got.GetTitle())
		assert.Equal(t, entities.Content("content"), got.GetContent())

		_, err = s.QueryByID(got.GetID())
		assert.NoError(t, err)

	})
}
