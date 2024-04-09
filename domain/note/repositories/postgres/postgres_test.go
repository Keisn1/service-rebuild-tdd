package postgres_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"database/sql"

	"github.com/Keisn1/note-taking-app/domain/note"
	"github.com/Keisn1/note-taking-app/domain/note/repositories/postgres"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

const (
	testDBName   = "test_note_taking_app"
	testUser     = "postgres"
	testPassword = "password"
)

func TestMain(m *testing.M) {
	exitCode := run(m)
	os.Exit(exitCode)
}

func TestNotesRepo_GetNoteByID(t *testing.T) {
	testDB, deleteTable := SetupNotesTable(t, fixtureNotes())
	defer testDB.Close()
	defer deleteTable()

	t.Run("Happy: get note by id", func(t *testing.T) {
		nR := postgres.NewNotesRepo(testDB)

		type testCase struct {
			noteID uuid.UUID
			want   note.Note
		}

		testCases := []testCase{
			{noteID: uuid.UUID{1}, want: note.MakeNote(uuid.UUID{1}, note.NewTitle("robs 1st note"), note.NewContent("robs 1st note content"), uuid.UUID{1})},
			{noteID: uuid.UUID{3}, want: note.MakeNote(uuid.UUID{3}, note.NewTitle("annas 1st note"), note.NewContent("annas 1st note content"), uuid.UUID{2})},
		}

		for _, tc := range testCases {
			got, err := nR.GetNoteByID(tc.noteID)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		}
	})

	t.Run("Note not found - sql.ErrNoRows", func(t *testing.T) {
		nR := postgres.NewNotesRepo(testDB)
		noteID := uuid.UUID{}

		_, err := nR.GetNoteByID(noteID)
		assert.EqualError(t, err, fmt.Errorf("getNoteByID: not found [%s]: %w", noteID, sql.ErrNoRows).Error())
	})
}

func TestNotesRepo_GetNotesByUserID(t *testing.T) {
	testDB, deleteTable := SetupNotesTable(t, fixtureNotes())
	defer testDB.Close()
	defer deleteTable()

	t.Run("Get notes by userID", func(t *testing.T) {
		nR := postgres.NewNotesRepo(testDB)
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
			got, err := nR.GetNotesByUserID(tc.userID)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.want, got)
		}
	})
	t.Run("Returns error on missing user", func(t *testing.T) {
		nR := postgres.NewNotesRepo(testDB)

		userID := uuid.UUID{}
		wantErrMsg := fmt.Sprintf("getNotesByUserID: not found [%s]", userID)
		_, err := nR.GetNotesByUserID(userID)
		assert.ErrorContains(t, err, wantErrMsg)
	})

	t.Run("Fowards error on database error", func(t *testing.T) {
		sdb := &SQLDB{}
		nR := postgres.NewNotesRepo(sdb)

		userID := uuid.UUID{}
		wantErr := fmt.Errorf("getNotesByUserID: [%s]: %w", userID, errors.New("DBError"))
		_, err := nR.GetNotesByUserID(userID)
		assert.EqualError(t, err, wantErr.Error())
	})

}
