package postgres

import (
	"database/sql"

	"github.com/Keisn1/note-taking-app/domain/note"
	"github.com/google/uuid"
)

type NoteRepo struct {
	db *sql.DB
}

func NewNotesRepo(db *sql.DB) NoteRepo {
	return NoteRepo{db: db}
}

type noteDB struct {
	id      uuid.UUID
	title   string
	content string
	userID  uuid.UUID
}

func (nR NoteRepo) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) {
	getNote := `
SELECT id, title, content, user_id FROM notes WHERE user_id=$1;
`
	rows, _ := nR.db.Query(getNote, userID)
	var notes []noteDB
	for rows.Next() {

		var n noteDB
		_ = rows.Scan(&n.id, &n.title, &n.content, &n.userID)
		notes = append(notes, n)
	}

	var ret []note.Note
	for _, n := range notes {
		ret = append(ret, note.MakeNote(
			n.id,
			note.NewTitle(n.title),
			note.NewContent(n.content),
			n.userID,
		))
	}

	return ret, nil
}
