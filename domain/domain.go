package domain

import "github.com/google/uuid"

var notes []Note

type NoteService struct{}

type NoteRepository interface {
	GetNoteByID(noteID uuid.UUID) Note
}

func (ns NoteService) GetNoteByID(noteID uuid.UUID) {
	if noteID == uuid.UUID([16]byte{1}) {
		return Note{
			ID:      noteID,
			Title:   "title1",
			Content: "content1",
		}
	}

	if noteID == uuid.UUID([16]byte{2}) {
		return Note{
			ID:      noteID,
			Title:   "title2",
			Content: "content2",
		}
	}
	return Note{}
}
