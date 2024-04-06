package usernoteSvc_test

import (
	"github.com/Keisn1/note-taking-app/services/usernoteSvc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("Get robs notes", func(t *testing.T) {
		want := usernoteSvc.Note{
			Owner:    "rob",
			NoteName: "robs note",
			NoteText: "robs note text",
		}

		got := usernoteSvc.GetNoteByName("rob")
		assert.Equal(t, want, got)
	})
}
