package postgres_test

import (
	"testing"

	"database/sql"

	"github.com/Keisn1/note-taking-app/domain/note"
	"github.com/Keisn1/note-taking-app/domain/note/repositories/postgres"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestPostGres(t *testing.T) {
	t.Run("Get notes of users by UserID", func(t *testing.T) {
		db := SetupPostGres(t, fixtureNotes())
		defer db.Close()

		notesR := postgres.NewNotesRepo(db)
		type testCase struct {
			userID uuid.UUID
			want   []note.Note
		}

		testCases := []testCase{
			{
				userID: uuid.UUID{1},
				want: []note.Note{
					note.MakeNote(uuid.UUID{1}, note.NewTitle("robs 1st note"), note.NewContent("robs 1st note content"), uuid.UUID{1}),
					note.MakeNote(uuid.UUID{2}, note.NewTitle("robs 2nd note"), note.NewContent("robs 2nd note content"), uuid.UUID{1}),
				},
			},
			{
				userID: uuid.UUID{2},
				want: []note.Note{
					note.MakeNote(uuid.UUID{3}, note.NewTitle("annas 1st note"), note.NewContent("annas 1st note content"), uuid.UUID{2}),
					note.MakeNote(uuid.UUID{4}, note.NewTitle("annas 2nd note"), note.NewContent("annas 2nd note content"), uuid.UUID{2}),
				},
			},
		}

		for _, tc := range testCases {
			got, err := notesR.GetNotesByUserID(tc.userID)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.want, got)
		}
	})

}

func fixtureNotes() []note.Note {
	return []note.Note{
		note.MakeNote(uuid.UUID{1}, note.NewTitle("robs 1st note"), note.NewContent("robs 1st note content"), uuid.UUID{1}),
		note.MakeNote(uuid.UUID{2}, note.NewTitle("robs 2nd note"), note.NewContent("robs 2nd note content"), uuid.UUID{1}),
		note.MakeNote(uuid.UUID{3}, note.NewTitle("annas 1st note"), note.NewContent("annas 1st note content"), uuid.UUID{2}),
		note.MakeNote(uuid.UUID{4}, note.NewTitle("annas 2nd note"), note.NewContent("annas 2nd note content"), uuid.UUID{2}),
	}
}

func SetupPostGres(t *testing.T, notes []note.Note) *sql.DB {
	url := "postgres://simba:mufassa@localhost:5432/test_note_taking_app"
	db, err := sql.Open("pgx", url)

	if err != nil {
		t.Fatal(err)
	}

	// dropDB := `DROP DATABASE IF EXISTS test_note_taking_app;`
	// createDB := `CREATE DATABASE test_note_taking_app;`
	dropNoteTable := `DROP TABLE IF EXISTS notes;`
	createNoteTable := `
CREATE TABLE notes(
id UUID PRIMARY KEY,
title TEXT,
content TEXT,
user_id UUID NOT NULL
)
`
	_, err = db.Exec(dropNoteTable)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(createNoteTable)
	if err != nil {
		t.Fatal(err)
	}

	insertRow := `INSERT INTO notes (id, title, content, user_id) VALUES ($1, $2, $3, $4)`
	for _, n := range notes {
		_, err = db.Exec(
			insertRow,
			n.GetID(),
			n.GetTitle().Get(),
			n.GetContent().Get(),
			n.GetUserID(),
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	return db
}

//TODO: need to have a timeout context for getting the database connection
