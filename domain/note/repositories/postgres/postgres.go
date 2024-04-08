package postgres

import (
	"database/sql"

	"github.com/Keisn1/note-taking-app/domain/note"
	"github.com/google/uuid"
)

type NoteRepo struct {
	conn *sql.Conn
}

func NewNotesRepo(conn *sql.Conn) NoteRepo {
	return NoteRepo{conn: conn}
}

func (nR NoteRepo) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) {
	return nil, nil
}
