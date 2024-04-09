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

func noteDBToNote(nDB noteDB) note.Note {
	return note.MakeNote(nDB.id,
		note.NewTitle(nDB.title),
		note.NewContent(nDB.content),
		nDB.userID)
}

func (nR NoteRepo) Create(n note.Note) error {
	return nil
}

func (nR NoteRepo) GetNoteByID(noteID uuid.UUID) (note.Note, error) {
	getNoteByID := `
	SELECT id, title, content, user_id FROM notes WHERE id=$1;
	`
	row := nR.db.QueryRow(getNoteByID, noteID)
	var nDB noteDB
	err := row.Scan(&nDB.id, &nDB.title, &nDB.content, &nDB.userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return note.Note{}, fmt.Errorf("getNoteByID: not found [%s]: %w", noteID, err)
		}
		return note.Note{}, fmt.Errorf("getNoteByID: [%s]: %w", noteID, err)
	}

	return noteDBToNote(nDB), nil
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
		var nDB noteDB
		err := rows.Scan(&nDB.id, &nDB.title, &nDB.content, &nDB.userID)
		if err != nil {
			return nil, fmt.Errorf("getNotesByUserId: [%s]: scan rows: %w", userID, err)
		}
		notes = append(notes, nDB)
	}

	if len(notes) == 0 {
		return nil, fmt.Errorf("getNotesByUserID: not found [%s]", userID)
	}

	var ret []note.Note
	for _, nDB := range notes {
		ret = append(ret, noteDBToNote(nDB))
	}

	return ret, nil
}
