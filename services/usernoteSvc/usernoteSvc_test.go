package usernoteSvc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("Get robs notes", func(t *testing.T) {
		want := struct {
			owner    string
			NoteName string
			NoteText string
		}{
			owner:    "rob",
			NoteName: "robs note",
			NoteText: "robs note text",
		}

		got := GetNoteByName("rob")
		assert.Equal(t, want, got)
	})
}
