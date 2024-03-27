package domain

type NotesStore interface {
	GetAllNotes() (Notes, error)
	GetNoteByUserIDAndNoteID(userID, noteID int) (Notes, error)
	GetNotesByUserID(userID int) (Notes, error)
	AddNote(userID int, note string) error
	EditNote(userID, noteID int, note string) error
	Delete(userID, noteID int) error
}
