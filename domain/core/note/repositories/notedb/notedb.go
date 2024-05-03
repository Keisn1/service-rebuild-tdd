package notedb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/google/uuid"
)

type dbNote struct {
	id      uuid.UUID
	title   string
	content string
	userID  uuid.UUID
}

type database interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

type NoteRepo struct {
	db database
}

func NewNotesRepo(db database) NoteRepo {
	return NoteRepo{db: db}
}

func (nR NoteRepo) Update(n note.Note) error {
	updateRow := `
	UPDATE notes
	SET title = $1, content = $2 WHERE id=$3 `

	res, err := nR.db.Exec(updateRow, n.Title.String(), n.Content.String(), n.ID)
	if err != nil {
		return fmt.Errorf("update: [%v]: %w", n, err)
	}
	if c, _ := res.RowsAffected(); c == 0 {
		return note.ErrNoteNotFound
	}

	return nil
}

func (nR NoteRepo) Delete(noteID uuid.UUID) error {
	deleteRow := `DELETE FROM notes WHERE id=$1`
	res, _ := nR.db.Exec(deleteRow, noteID)

	if c, _ := res.RowsAffected(); c == 0 {
		return note.ErrNoteNotFound
	}

	return nil
}

func (nR NoteRepo) Create(n note.Note) error {
	insertRow := `INSERT INTO notes (id, title, content, user_id) VALUES ($1, $2, $3, $4)`
	_, err := nR.db.Exec(
		insertRow,
		n.ID,
		n.Title.String(),
		n.Content.String(),
		n.UserID,
	)
	if err != nil {
		return fmt.Errorf("create: [%s]", n.ID)
	}

	return nil
}

func (nR NoteRepo) QueryByID(ctx context.Context, noteID uuid.UUID) (note.Note, error) {
	queryByIDSqlStmt := `
	SELECT id, title, content, user_id FROM notes WHERE id=$1;
	`
	row := nR.db.QueryRowContext(ctx, queryByIDSqlStmt, noteID)
	var nDB dbNote
	err := row.Scan(&nDB.id, &nDB.title, &nDB.content, &nDB.userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return note.Note{}, note.ErrNoteNotFound
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

	var notes []dbNote
	for rows.Next() {
		var nDB dbNote
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

func noteDBToNote(nDB dbNote) note.Note {
	return note.Note{
		ID:      nDB.id,
		Title:   note.NewTitle(nDB.title),
		Content: note.NewContent(nDB.content),
		UserID:  nDB.userID,
	}
}
