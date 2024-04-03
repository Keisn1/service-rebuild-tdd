package domain

import "github.com/google/uuid"

type NotesStore interface {
	GetAllNotes() (Notes, error)
	GetNoteByUserIDAndNoteID(userID uuid.UUID, noteID int) (Notes, error)
	GetNotesByUserID(userID uuid.UUID) (Notes, error)
	AddNote(userID uuid.UUID, np NotePost) error
	EditNote(userID, noteID int, note string) error
	Delete(userID uuid.UUID, noteID int) error
}
