package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Keisn1/note-taking-app/domain/note"
	"github.com/google/uuid"
)

type database interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type NoteRepo struct {
	db database
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

func (nR NoteRepo) GetNoteByID(noteID uuid.UUID) (note.Note, error) {
	getNoteByID := `
	SELECT id, title, content, user_id FROM notes WHERE id=$1;
	`
	row := nR.db.QueryRow(getNoteByID, noteID)
	var nDB noteDB
	_ = row.Scan(&nDB.id, &nDB.title, &nDB.content, &nDB.userID)

	n := note.MakeNote(nDB.id,
		note.NewTitle(nDB.title),
		note.NewContent(nDB.content),
		nDB.userID)

	return n, nil
}

func (nR NoteRepo) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) {
	getNotesByUserID := `
	SELECT id, title, content, user_id FROM notes WHERE user_id=$1;
	`
	rows, err := nR.db.Query(getNotesByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("getNotesByUserID: [%s]: %w", userID, err)
	}
	defer rows.Close()

	var notes []noteDB
	for rows.Next() {
		var n noteDB
		err := rows.Scan(&n.id, &n.title, &n.content, &n.userID)
		if err != nil {
			return nil, fmt.Errorf("getNotesByUserId: [%s]: scan rows: %w", userID, err)
		}
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
