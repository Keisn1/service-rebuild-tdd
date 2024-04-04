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

	if nID == uuid.UUID([16]byte{1}) {
		return usernote.NewUserNote(nID, "", "", uuid.UUID([16]byte{2}))
	}

	if nID == uuid.UUID([16]byte{3}) {
		return usernote.NewUserNote(nID, "", "", uuid.UUID([16]byte{4}))
	}
	return usernote.UserNote{}
	// if noteID == uuid.UUID([16]byte{2}) {
	// 	return Note{
	// 		ID:      noteID,
	// 		Title:   "title2",
	// 		Content: "content2",
	// 	}
	// }
	// return Note{}
}
