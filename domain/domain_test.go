package domain_test

import (
	"testing"

	"github.com/Keisn1/note-taking-app/domain"
	"github.com/Keisn1/note-taking-app/domain/entities"
	"github.com/Keisn1/note-taking-app/domain/usernote"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNotes(t *testing.T) {
	// t.Run("When I add a note, I can get this note in return by its ID", func(t *testing.T) {
	// 	notes := []domain.Note{}
	// 	wantNote := domain.NewNote("title1", "content1")

	// 	domain.AddNote(wantNote)

	// 	gotNote := domain.GetNoteByID(wantNote.ID)

	// 	assert.Equal(t, wantNote, gotNote)
	// })
	t.Run("Return note for noteID", func(t *testing.T) {
		nID := uuid.UUID([16]byte{1})
		uID := uuid.UUID([16]byte{2})
		want := usernote.NewUserNote(nID, "title1", "content1", uID)

		got := domain.GetNoteByID(nID)
		assert.Equal(t, want, got)

		// noteID = uuid.UUID([16]byte{2})
		// want = domain.Note{
		// 	ID:      noteID,
		// 	Title:   "title2",
		// 	Content: "content2",
		// }
		// got = domain.GetNoteByID(noteID)
		// assert.Equal(t, want, got)
	})
}
