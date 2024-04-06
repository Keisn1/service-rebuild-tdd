package note_test

import (
	"testing"

	"github.com/Keisn1/note-taking-app/domain/note"
	"github.com/stretchr/testify/assert"
)

func TestNote(t *testing.T) {
	t.Run("I can set a title and get a title", func(t *testing.T) {
		t := note.NewTitle("title")
		t.Set("newTitle")

		want := "newTitle"
		got := t.Get()
		assert.Equal(t, "newTitle", got)
	})
}
