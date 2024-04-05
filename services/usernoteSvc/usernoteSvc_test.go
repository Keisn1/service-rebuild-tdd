package usernoteSvc_test

import (
	"fmt"
	"testing"

	"github.com/Keisn1/note-taking-app/domain/entities"
	"github.com/Keisn1/note-taking-app/domain/user"
	"github.com/Keisn1/note-taking-app/domain/usernote"
	"github.com/Keisn1/note-taking-app/services/usernoteSvc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	usr1 := user.NewUser("", "")
	usr2 := user.NewUser("", "")
	userID1 := usr1.GetID()
	userID2 := usr2.GetID()

	note1 := usernote.NewUserNote("", "", userID1)
	note2 := usernote.NewUserNote("", "", userID1)
	note3 := usernote.NewUserNote("", "", userID2)
	note4 := usernote.NewUserNote("", "", userID2)

	noteID1 := note1.GetID()
	noteID2 := note2.GetID()
	noteID3 := note3.GetID()
	noteID4 := note4.GetID()

	un := &StubUserNoteRepository{
		Usernotes: map[uuid.UUID]usernote.UserNote{
			noteID1: note1, noteID2: note2, noteID3: note3, noteID4: note4,
		},
	}

	u := &StubUserRepository{users: map[uuid.UUID]user.User{
		usr1.GetID(): usr1,
		usr2.GetID(): usr2,
	}}

	s := usernoteSvc.NewUserNoteService(usernoteSvc.WithUserNoteRepository(un), usernoteSvc.WithUserRepository(u))

	t.Run("Get note by noteID", func(t *testing.T) {
		noteID := noteID1
		want := note1
		got, err := s.QueryByID(noteID)
		assert.NoError(t, err)
		assert.Equal(t, want, got)

		noteID = noteID2
		want = note2
		got, err = s.QueryByID(noteID)
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
		userID := userID1
		wantNotes := []usernote.UserNote{note1, note2}
		gotNotes, err := s.QueryByUserID(userID)
		assert.NoError(t, err)
		assert.ElementsMatch(t, wantNotes, gotNotes)

		userID = userID2
		wantNotes = []usernote.UserNote{note3, note4}
		gotNotes, err = s.QueryByUserID(userID)
		assert.NoError(t, err)
		assert.ElementsMatch(t, wantNotes, gotNotes)
	})

	t.Run("Return error if no notes found for userID", func(t *testing.T) {
		userIDx := uuid.UUID([16]byte{100})
		_, err := s.QueryByUserID(userIDx)
		expectedErrorSubString := fmt.Sprintf("querybyuserid: userID[%s]", userIDx)
		assert.ErrorContains(t, err, expectedErrorSubString)
	})

	t.Run("Create a note", func(t *testing.T) {
		userID := userID1
		got, err := s.Create(userID, "title", "content")
		assert.NoError(t, err)
		assert.Equal(t, got.GetUserID(), userID1)
		assert.Equal(t, entities.Title("title"), got.GetTitle())
		assert.Equal(t, entities.Content("content"), got.GetContent())

		want := got
		got, err = s.QueryByID(got.GetID())
		assert.NoError(t, err)
		assert.Equal(t, got, want)
	})

	t.Run("When User can not be retrieved, then the note can not be created", func(t *testing.T) {
		userIDx := uuid.UUID([16]byte{100})
		_, err := s.Create(userIDx, "title", "content")
		expectedErrorSubString := fmt.Sprintf("Create: userID[%s]", userIDx)
		assert.ErrorContains(t, err, expectedErrorSubString)
	})

	t.Run("Edit a title of a note", func(t *testing.T) {
		noteID := noteID1

		un.EditWasCalled = false
		note, err := s.QueryByID(noteID)
		assert.NoError(t, err)
		formerTitle := note.GetTitle()
		formerContent := note.GetContent()

		got, err := s.Edit(noteID, "title", "content")
		assert.NoError(t, err)
		assert.Equal(t, got.GetID(), noteID)
		assert.Equal(t, got.GetTitle(), entities.Title("title"))
		assert.Equal(t, got.GetContent(), entities.Content("content"))

		if !un.EditWasCalled {
			got, err := s.QueryByID(noteID)
			assert.NoError(t, err)
			assert.Equal(t, got.GetTitle(), formerTitle)
			assert.Equal(t, got.GetContent(), formerContent)
		}

		// got2, err := s.QueryByID(noteID)
		// assert.NoError(t, err)
		// assert.NotEqual(t, &got1, &got2)
		// assert.Equal(t, got.GetID(), noteID)
		// assert.Equal(t, got.GetTitle(), entities.Title("title"))
		// assert.Equal(t, got.GetContent(), entities.Content("content"))
	})

}
