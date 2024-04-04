package domain_test

import (
	"testing"

	"github.com/Keisn1/note-taking-app/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetNotes(t *testing.T) {
	wantNotes := []string{"note1", "note2"}

	gotNotes := domain.GetNotes()
	assert.Equal(t, wantNotes, gotNotes)

	t.Run("When i add a note, i can get this note in return by its ID", func(t *testing.T) {
		wantNote := NewNote("note1")

		domain.AddNote(wantNote)

		gotNote := domain.GetNoteByID(wantNote.noteID)

		assert.Equal(t, wantNote, gotNote)
	})
}
