package usernoteSvc_test

import (
	"fmt"
	"testing"

	"github.com/Keisn1/note-taking-app/domain/entities"
	"github.com/Keisn1/note-taking-app/domain/usernote"
	"github.com/Keisn1/note-taking-app/services/usernoteSvc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type StubUserNoteRepository struct {
	usernotes map[uuid.UUID]usernote.UserNote
}

func (sUNR *StubUserNoteRepository) Create(userID uuid.UUID, title, content string) (usernote.UserNote, error) {
	u := usernote.NewUserNote(title, content, userID)
	sUNR.usernotes[u.GetID()] = u
	return u, nil
}

func (sUNR *StubUserNoteRepository) GetNoteByID(noteID uuid.UUID) (usernote.UserNote, error) {
	n, ok := sUNR.usernotes[noteID]
	if !ok {
		return usernote.UserNote{}, fmt.Errorf("Note note found")
	}
	return n, nil
}

func (sUNR *StubUserNoteRepository) GetNotesByUserID(userID uuid.UUID) ([]usernote.UserNote, error) {
	var ret []usernote.UserNote
	for _, n := range sUNR.usernotes {
		if n.GetUserID() == userID {
			ret = append(ret, n)
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("No note found for userID[%s]", userID)
	}
	return ret, nil
}

func TestService(t *testing.T) {
	userID1 := uuid.UUID([16]byte{1})
	userID2 := uuid.UUID([16]byte{2})

	note1 := usernote.NewUserNote("", "", userID1)
	note2 := usernote.NewUserNote("", "", userID1)
	note3 := usernote.NewUserNote("", "", userID2)
	note4 := usernote.NewUserNote("", "", userID2)

	noteID1 := note1.GetID()
	noteID2 := note2.GetID()
	noteID3 := note3.GetID()
	noteID4 := note4.GetID()

	u := &StubUserNoteRepository{
		usernotes: map[uuid.UUID]usernote.UserNote{
			noteID1: note1, noteID2: note2, noteID3: note3, noteID4: note4,
		},
	}
	s := usernoteSvc.NewUserNoteService(usernoteSvc.WithUserNoteRepository(u))

	t.Run("Get note by noteID", func(t *testing.T) {
		want := note1
		got, err := s.QueryByID(noteID1)
		assert.NoError(t, err)
		assert.Equal(t, want, got)

		want = note2
		got, err = s.QueryByID(noteID2)
		assert.Equal(t, want, got)
		assert.NoError(t, err)
	})

	t.Run("Return error for missing note", func(t *testing.T) {
		noteIDx := uuid.UUID([16]byte{100})
		_, err := s.QueryByID(noteIDx)
		expectedErrorSubString := fmt.Sprintf("querybyid: noteID[%s]", noteIDx)
		assert.ErrorContains(t, err, expectedErrorSubString)
	})

	t.Run("Return notes by UserID", func(t *testing.T) {
		wantNotes := []usernote.UserNote{note1, note2}
		gotNotes, err := s.QueryByUserID(userID1)
		assert.NoError(t, err)
		assert.ElementsMatch(t, wantNotes, gotNotes)

		wantNotes = []usernote.UserNote{note3, note4}
		gotNotes, err = s.QueryByUserID(userID2)
		assert.NoError(t, err)
		assert.ElementsMatch(t, wantNotes, gotNotes)
	})

	t.Run("Return error if no notes found for userID", func(t *testing.T) {
		userIDx := uuid.UUID([16]byte{100})
		_, err := s.QueryByUserID(userIDx)
		expectedErrorSubString := fmt.Sprintf("querybyuserid: userID[%s]", userIDx)
		assert.ErrorContains(t, err, expectedErrorSubString)
	})

	t.Run("Add a note", func(t *testing.T) {
		userID := uuid.UUID([16]byte{1})
		got, err := s.Create(userID, "title", "content")
		assert.NoError(t, err)
		assert.Equal(t, got.GetUserID(), userID)
		assert.Equal(t, entities.Title("title"), got.GetTitle())
		assert.Equal(t, entities.Content("content"), got.GetContent())

		want := got
		got, err = s.QueryByID(got.GetID())
		assert.NoError(t, err)
		assert.Equal(t, got, want)
	})

	t.Run("Add a note", func(t *testing.T) {
		userID := uuid.UUID([16]byte{1})
		got, err := s.Create(userID, "title", "content")
		assert.NoError(t, err)
		assert.Equal(t, got.GetUserID(), userID)
		assert.Equal(t, entities.Title("title"), got.GetTitle())
		assert.Equal(t, entities.Content("content"), got.GetContent())

		want := got
		got, err = s.QueryByID(got.GetID())
		assert.NoError(t, err)
		assert.Equal(t, got, want)
	})

	t.Run("When User can not be retrieved, then Add throws error", func(t *testing.T) {
		userIDx := uuid.UUID([16]byte{100})
		_, err := s.Create(userIDx, "title", "content")
		expectedErrorSubString := fmt.Sprintf("Create: userID[%s]", userIDx)
		assert.ErrorContains(t, err, expectedErrorSubString)
	})
}
