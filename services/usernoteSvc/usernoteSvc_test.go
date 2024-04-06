package usernoteSvc_test

import (
	"testing"

	"github.com/Keisn1/note-taking-app/services/usernoteSvc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("Get notes by UserID", func(t *testing.T) {
		want := []usernoteSvc.Note{
			{
				UserID:  uuid.UUID{1},
				Title:   "robs 1st note",
				Content: "robs 1st note content",
			},
			{
				UserID:  uuid.UUID{1},
				Title:   "robs 2nd note",
				Content: "robs 2nd note content",
			},
		}

		got := usernoteSvc.GetNotesByUserID(uuid.UUID{1})
		assert.ElementsMatch(t, want, got)
	})
}
