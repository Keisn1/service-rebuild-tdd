package notedb_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/Keisn1/note-taking-app/domain/core/note/repositories/notedb"
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

func TestNotesRepo_Update(t *testing.T) {
	testDB, deleteTable := SetupNotesTable(t, fixtureNotes())
	defer testDB.Close()
	defer deleteTable()

	t.Run("Given an error received by the DB, the error is forwarded", func(t *testing.T) {
		nR := notedb.NewNotesRepo(&stubSQLDB{})
		n := note.Note{ID: uuid.New(), Title: note.NewTitle(""), Content: note.NewContent(""), UserID: uuid.New()}
		err := nR.Update(n)
		assert.ErrorContains(t, err, fmt.Sprintf("update: [%v]: DBError", n))
	})

	t.Run("Given a note NOT present in the system, return ErrNoteNotFound", func(t *testing.T) {
		nR := notedb.NewNotesRepo(testDB)
		n := note.Note{ID: uuid.New(), Title: note.NewTitle(""), Content: note.NewContent(""), UserID: uuid.New()}
		err := nR.Update(n)
		assert.ErrorContains(t, err, note.ErrNoteNotFound.Error())
	})

	t.Run("Given a note present in the system, I can update its title and its content", func(t *testing.T) {
		nR := notedb.NewNotesRepo(testDB)
		n := note.Note{ID: uuid.UUID{1}, Title: note.NewTitle("new title"), Content: note.NewContent("new content"), UserID: uuid.UUID{1}}

		err := nR.Update(n)
		assert.NoError(t, err)

		got, err := nR.QueryByID(context.Background(), n.ID)
		assert.NoError(t, err)
		assert.Equal(t, n, got)
	})
}

func TestNotesRepo_Delete(t *testing.T) {
	testDB, deleteTable := SetupNotesTable(t, fixtureNotes())
	defer testDB.Close()
	defer deleteTable()

	t.Run("Able to delete a note", func(t *testing.T) {
		nR := notedb.NewNotesRepo(testDB)
		ctx := context.Background()
		noteID := uuid.New()
		n := note.Note{ID: noteID, Title: note.NewTitle("new title"), Content: note.NewContent("new content"), UserID: uuid.UUID{1}}

		nR.Create(n)
		got, err := nR.QueryByID(ctx, noteID)
		assert.NoError(t, err)
		assert.Equal(t, n, got)

		err = nR.Delete(noteID)
		assert.NoError(t, err)

		_, err = nR.QueryByID(ctx, noteID)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "not found")
	})
	t.Run("Delete non-present note throws error ", func(t *testing.T) {
		nR := notedb.NewNotesRepo(testDB)

		noteID := uuid.New()
		err := nR.Delete(noteID)
		assert.ErrorContains(t, err, note.ErrNoteNotFound.Error())
		assert.ErrorContains(t, err, "not found")
	})
}

func TestNotesRepo_Create(t *testing.T) {
	t.Run("Add a note", func(t *testing.T) {
		testDB, deleteTable := SetupNotesTable(t, []note.Note{})
		defer testDB.Close()
		defer deleteTable()
		nR := notedb.NewNotesRepo(testDB)

		ctx := context.Background()
		n := note.Note{ID: uuid.UUID{1}, Title: note.NewTitle("new title"), Content: note.NewContent("new content"), UserID: uuid.UUID{1}}

		err := nR.Create(n)
		assert.NoError(t, err)

		got, err := nR.QueryByID(ctx, n.ID)
		assert.NoError(t, err)
		assert.Equal(t, got, n)
	})

	t.Run("Throws error if note to be created already exists", func(t *testing.T) {
		testDB, deleteTable := SetupNotesTable(t, fixtureNotes())
		defer testDB.Close()
		defer deleteTable()
		nR := notedb.NewNotesRepo(testDB)

		n := note.Note{ID: uuid.UUID{1}, Title: note.NewTitle("new title"), Content: note.NewContent("new content"), UserID: uuid.UUID{1}}
		err := nR.Create(n)
		assert.Error(t, err)
		assert.ErrorContains(t, err, fmt.Sprintf("create: [%s]", n.ID))
	})
}

func TestNotesRepo_QueryByID(t *testing.T) {
	testDB, deleteTable := SetupNotesTable(t, fixtureNotes())
	defer testDB.Close()
	defer deleteTable()

	t.Run("Happy: get note by id", func(t *testing.T) {
		nR := notedb.NewNotesRepo(testDB)

		type testCase struct {
			noteID uuid.UUID
			want   note.Note
		}

		testCases := []testCase{
			{
				noteID: uuid.UUID{1},
				want: note.Note{
					ID: uuid.UUID{1}, Title: note.NewTitle("robs 1st note"), Content: note.NewContent("robs 1st note content"), UserID: uuid.UUID{1},
				},
			},
			{
				noteID: uuid.UUID{3},
				want: note.Note{
					ID: uuid.UUID{3}, Title: note.NewTitle("annas 1st note"), Content: note.NewContent("annas 1st note content"), UserID: uuid.UUID{2},
				},
			},
		}

		for _, tc := range testCases {
			got, err := nR.QueryByID(context.Background(), tc.noteID)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		}
	})

	t.Run("Note not found", func(t *testing.T) {
		ctx := context.Background()
		nR := notedb.NewNotesRepo(testDB)
		noteID := uuid.UUID{}

		_, err := nR.QueryByID(ctx, noteID)
		assert.EqualError(t, err, note.ErrNoteNotFound.Error())
	})
}

func TestNotesRepo_GetNotesByUserID(t *testing.T) {
	testDB, deleteTable := SetupNotesTable(t, fixtureNotes())
	defer testDB.Close()
	defer deleteTable()

	t.Run("Get notes by userID", func(t *testing.T) {
		nR := notedb.NewNotesRepo(testDB)
		type testCase struct {
			userID uuid.UUID
			want   []note.Note
		}

		testCases := []testCase{
			{
				userID: uuid.UUID{1},
				want: []note.Note{
					{ID: uuid.UUID{1}, Title: note.NewTitle("robs 1st note"), Content: note.NewContent("robs 1st note content"), UserID: uuid.UUID{1}},
					{ID: uuid.UUID{2}, Title: note.NewTitle("robs 2nd note"), Content: note.NewContent("robs 2nd note content"), UserID: uuid.UUID{1}},
				},
			},
			{
				userID: uuid.UUID{2},
				want: []note.Note{
					{ID: uuid.UUID{3}, Title: note.NewTitle("annas 1st note"), Content: note.NewContent("annas 1st note content"), UserID: uuid.UUID{2}},
					{ID: uuid.UUID{4}, Title: note.NewTitle("annas 2nd note"), Content: note.NewContent("annas 2nd note content"), UserID: uuid.UUID{2}},
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
		nR := notedb.NewNotesRepo(testDB)

		userID := uuid.UUID{}
		wantErrMsg := fmt.Sprintf("getNotesByUserID: not found [%s]", userID)
		_, err := nR.GetNotesByUserID(userID)
		assert.ErrorContains(t, err, wantErrMsg)
	})

	t.Run("Fowards error on database error", func(t *testing.T) {
		stubDB := &stubSQLDB{}
		nR := notedb.NewNotesRepo(stubDB)

		userID := uuid.UUID{}
		wantErr := fmt.Errorf("getNotesByUserID: [%s]: %w", userID, errors.New("DBError"))
		_, err := nR.GetNotesByUserID(userID)
		assert.EqualError(t, err, wantErr.Error())
	})

}
