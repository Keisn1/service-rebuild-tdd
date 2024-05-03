package note_test

import (
	"testing"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNote(t *testing.T) {
	t.Run("An empty title returns string-zeroValue for String()", func(t *testing.T) {
		tt := note.Title{}
		assert.Equal(t, "", tt.String())
	})

	t.Run("An empty content returns string-zeroValue for String()", func(t *testing.T) {
		c := note.Content{}
		assert.Equal(t, "", c.String())
	})

	t.Run("I can create a new note and access their fields via getter and setter", func(t *testing.T) {
		noteID, userID := uuid.New(), uuid.New()
		n := note.Note{NoteID: noteID, Title: note.NewTitle("title"), Content: note.NewContent("content"), UserID: userID}

		assert.Equal(t, noteID, n.GetID())
		assert.Equal(t, note.NewTitle("title"), n.GetTitle())
		assert.Equal(t, note.NewContent("content"), n.GetContent())
		assert.Equal(t, userID, n.GetUserID())

		newNoteID := uuid.New()
		n.SetID(newNoteID)
		assert.Equal(t, newNoteID, n.GetID())

		n.SetTitle("new title")
		assert.Equal(t, "new title", n.GetTitle().String())

		n.SetContent("new content")
		assert.Equal(t, "new content", n.GetContent().String())
	})

	t.Run("I can set a title and get a title", func(t *testing.T) {
		title := note.NewTitle("title")
		title.Set("newTitle")

		want := "newTitle"
		got := title.String()
		assert.Equal(t, want, got)
	})

	t.Run("I can check if title is empty", func(t *testing.T) {
		title := note.Title{}
		assert.True(t, title.IsEmpty())
	})

	t.Run("I can set a content and get a content", func(t *testing.T) {
		content := note.NewContent("content")
		content.Set("newContent")

		want := "newContent"
		got := content.String()
		assert.Equal(t, want, got)
	})

	t.Run("I can check if content is empty", func(t *testing.T) {
		content := note.Content{}
		assert.True(t, content.IsEmpty())
	})
}
