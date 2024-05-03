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
		n := note.NewNote(uuid.New(), "", "", uuid.New())
		err := nR.Update(n)
		assert.ErrorContains(t, err, fmt.Sprintf("update: [%v]: DBError", n))
	})

	t.Run("Given a note NOT present in the system, return ErrNoteNotFound", func(t *testing.T) {
		nR := notedb.NewNotesRepo(testDB)

		n := note.NewNote(uuid.New(), "", "", uuid.New())
		err := nR.Update(n)
		assert.ErrorContains(t, err, note.ErrNoteNotFound.Error())
	})

	t.Run("Given a note present in the system, I can update its title and its content", func(t *testing.T) {
		nR := notedb.NewNotesRepo(testDB)

		type testCase struct {
			name string
			ctx  context.Context
			n    note.Note
			want note.Note
		}

		testCases := []testCase{
			{
				name: "New Title, 0length content, update of both: 'new title' and '' ",
				n:    note.NewNote(uuid.UUID{1}, "new title", "new content", uuid.UUID{1}),
				want: note.NewNote(uuid.UUID{1}, "new title", "new content", uuid.UUID{1}),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := nR.Update(tc.n)
				assert.NoError(t, err)

				got, err := nR.QueryByID(tc.ctx, tc.n.GetID())
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
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
		nN := note.NewNote(noteID, "title", "content", uuid.UUID{1})

		nR.Create(nN)
		got, err := nR.QueryByID(ctx, noteID)
		assert.NoError(t, err)
		assert.Equal(t, nN, got)

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
		n := note.NewNote(uuid.UUID{1}, "robs 1st note", "robs 1st note content", uuid.UUID{1})
		err := nR.Create(n)
		assert.NoError(t, err)

		got, err := nR.QueryByID(ctx, n.GetID())
		assert.NoError(t, err)
		assert.Equal(t, got, n)
	})

	t.Run("Throws error if note to be created already exists", func(t *testing.T) {
		testDB, deleteTable := SetupNotesTable(t, fixtureNotes())
		defer testDB.Close()
		defer deleteTable()
		nR := notedb.NewNotesRepo(testDB)

		n := note.NewNote(uuid.UUID{1}, "robs 1st note", "robs 1st note content", uuid.UUID{1})
		err := nR.Create(n)
		assert.Error(t, err)
		assert.ErrorContains(t, err, fmt.Sprintf("create: [%s]", n.GetID()))
	})
}

func TestNotesRepo_QueryByID(t *testing.T) {
	testDB, deleteTable := SetupNotesTable(t, fixtureNotes())
	defer testDB.Close()
	defer deleteTable()

	t.Run("Happy: get note by id", func(t *testing.T) {
		nR := notedb.NewNotesRepo(testDB)

		type testCase struct {
			ctx    context.Context
			noteID uuid.UUID
			want   note.Note
		}

		testCases := []testCase{
			{noteID: uuid.UUID{1}, want: note.NewNote(uuid.UUID{1}, "robs 1st note", "robs 1st note content", uuid.UUID{1})},
			{noteID: uuid.UUID{3}, want: note.NewNote(uuid.UUID{3}, "annas 1st note", "annas 1st note content", uuid.UUID{2})},
		}

		for _, tc := range testCases {
			got, err := nR.QueryByID(tc.ctx, tc.noteID)
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

	t.Run("Test for context timeout", func(t *testing.T) {
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
					note.NewNote(uuid.UUID{1}, "robs 1st note", "robs 1st note content", uuid.UUID{1}),
					note.NewNote(uuid.UUID{2}, "robs 2nd note", "robs 2nd note content", uuid.UUID{1}),
				},
			},
			{
				userID: uuid.UUID{2},
				want: []note.Note{
					note.NewNote(uuid.UUID{3}, "annas 1st note", "annas 1st note content", uuid.UUID{2}),
					note.NewNote(uuid.UUID{4}, "annas 2nd note", "annas 2nd note content", uuid.UUID{2}),
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
