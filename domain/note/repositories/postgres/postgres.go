package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Keisn1/note-taking-app/domain/note"
	"github.com/google/uuid"
)

type NoteRepo struct {
	db database
}

type database interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func NewNotesRepo(db database) NoteRepo {
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
	defer rows.Close()

	var notes []noteDB
	for rows.Next() {
		var n noteDB
		_ = rows.Scan(&n.id, &n.title, &n.content, &n.userID)
		notes = append(notes, n)
	}

	if len(notes) == 0 {
		return nil, fmt.Errorf("getNotesByUserID: not found [%s]", userID)
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
