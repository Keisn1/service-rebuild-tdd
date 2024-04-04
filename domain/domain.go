package domain

import (
	"github.com/Keisn1/note-taking-app/domain/usernote"
	"github.com/google/uuid"
)

type NoteService struct{}

// type NoteRepository interface {
// 	GetNoteByID(noteID uuid.UUID) Note
// }

func (ns NoteService) GetNoteByID(nID uuid.UUID) usernote.UserNote {
	return usernote.NewUserNote(nID, "title1", "content1", uuid.UUID([16]byte{2}))
	// if noteID == uuid.UUID([16]byte{1}) {
	// 	return Note{
	// 		ID:      noteID,
	// 		Title:   "title1",
	// 		Content: "content1",
	// 	}
	// }

	// if noteID == uuid.UUID([16]byte{2}) {
	// 	return Note{
	// 		ID:      noteID,
	// 		Title:   "title2",
	// 		Content: "content2",
	// 	}
	// }
	// return Note{}
}
