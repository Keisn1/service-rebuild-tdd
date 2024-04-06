package usernoteSvc_test

import (
	"testing"

	"github.com/Keisn1/note-taking-app/domain/usernote"
	"github.com/Keisn1/note-taking-app/services/usernoteSvc"
	"github.com/google/uuid"
)

func TestService(t *testing.T) {
	t.Run("Get note by noteID", func(t *testing.T) {
		un1 := UserNoteDB{
			ID:      uuid.New(),
			Title:   "title1",
			Content: "content1",
			UserID:  uuid.New(),
		}

		un2 := UserNoteDB{
			ID:      uuid.New(),
			Title:   "title2",
			Content: "content2",
			UserID:  uuid.New(),
		}

		sUNR := &StubUserNoteRepository{
			Usernotes: map[uuid.UUID]UserNoteDB{
				un1.ID: un1,
				un2.ID: un2,
			},
		}
		s := usernoteSvc.NewUserNoteService(usernoteSvc.WithUserNoteRepository(sUNR))

		_, err := s.QueryByID()
		assert.NoError(t, err)
	})

	// t.Run("Edits on retrieved note don't have effect on data in store", func(t *testing.T) {
	// 	got, err := s.QueryByID(noteID1)
	// 	assert.NoError(t, err)

	// 	// title := got.GetTitle()
	// 	// content := got.GetContent()

	// 	got.SetTitle("new title")
	// 	got.SetContent("new content")

	// 	got2, err := s.QueryByID(noteID1)
	// 	assert.NoError(t, err)
	// 	assert.NotEqual(t, got, got2)
	// })

	// t.Run("Return error for missing note", func(t *testing.T) {
	// 	noteIDx := uuid.UUID([16]byte{100})
	// 	_, err := s.QueryByID(noteIDx)
	// 	expectedErrorSubString := fmt.Sprintf("querybyid: noteID[%s]", noteIDx)
	// 	assert.ErrorContains(t, err, expectedErrorSubString)
	// })

	// t.Run("Return notes by UserID", func(t *testing.T) {
	// 	userID := userID1
	// 	wantNotes := []usernote.UserNote{note1, note2}
	// 	gotNotes, err := s.QueryByUserID(userID)
	// 	assert.NoError(t, err)
	// 	assert.ElementsMatch(t, wantNotes, gotNotes)

	// 	userID = userID2
	// 	wantNotes = []usernote.UserNote{note3, note4}
	// 	gotNotes, err = s.QueryByUserID(userID)
	// 	assert.NoError(t, err)
	// 	assert.ElementsMatch(t, wantNotes, gotNotes)
	// })

	// t.Run("Return error if no notes found for userID", func(t *testing.T) {
	// 	userIDx := uuid.UUID([16]byte{100})
	// 	_, err := s.QueryByUserID(userIDx)
	// 	expectedErrorSubString := fmt.Sprintf("querybyuserid: userID[%s]", userIDx)
	// 	assert.ErrorContains(t, err, expectedErrorSubString)
	// })

	// t.Run("Create a note", func(t *testing.T) {
	// 	userID := userID1
	// 	got, err := s.Create(userID, "title", "content")
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, got.GetUserID(), userID1)
	// 	assert.Equal(t, entities.Title("title"), got.GetTitle())
	// 	assert.Equal(t, entities.Content("content"), got.GetContent())

	// 	want := got
	// 	got, err = s.QueryByID(got.GetID())
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, got, want)
	// })

	// t.Run("When User can not be retrieved, then the note can not be created", func(t *testing.T) {
	// 	userIDx := uuid.UUID([16]byte{100})
	// 	_, err := s.Create(userIDx, "title", "content")
	// 	expectedErrorSubString := fmt.Sprintf("Create: userID[%s]", userIDx)
	// 	assert.ErrorContains(t, err, expectedErrorSubString)
	// })

	// t.Run("Edit a title of a note", func(t *testing.T) {
	// 	noteID := noteID1

	// 	got1, err := s.Edit(noteID, "title", "")
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, got1.GetID(), noteID)
	// 	assert.Equal(t, got1.GetTitle(), entities.Title("title"))

	// 	got2, err := s.QueryByID(noteID)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, got1, got2)
	// 	assert.Equal(t, got2.GetID(), noteID)
	// 	assert.Equal(t, got2.GetTitle(), entities.Title("title"))
	// })
}
