package domain

import "github.com/google/uuid"

type NotesStore interface {
	GetAllNotes() (Notes, error)
	GetNoteByUserIDAndNoteID(userID, noteID int) (Notes, error)
	GetNotesByUserID(userID int) (Notes, error)
	AddNote(userID uuid.UUID, np NotePost) error
	EditNote(userID, noteID int, note string) error
	Delete(userID, noteID int) error
}
