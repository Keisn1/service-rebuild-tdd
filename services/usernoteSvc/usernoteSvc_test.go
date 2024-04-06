package usernoteSvc_test

import (
	"testing"

	"github.com/Keisn1/note-taking-app/services/usernoteSvc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("Get note by UserID", func(t *testing.T) {
		want := usernoteSvc.Note{
			UserID:  uuid.UUID{1},
			Title:   "robs note",
			Content: "robs note content",
		}

		got := usernoteSvc.GetNoteByUserID(uuid.UUID{1})
		assert.Equal(t, want, got)
	})
}
