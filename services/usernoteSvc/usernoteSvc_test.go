package usernoteSvc_test

import (
	"testing"

	"github.com/Keisn1/note-taking-app/domain/note"
	svc "github.com/Keisn1/note-taking-app/services/usernoteSvc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNotes(t *testing.T) {
	t.Run("Given a note not present in the system, return error", func(t *testing.T) {
		notesR, err := note.NewNotesRepo(fixtureNotes())
		assert.NoError(t, err)
		notesS := svc.NewNotesService(notesR)

		noteID := uuid.New()
		_, err = notesS.Update(noteID, "some title")
		assert.ErrorContains(t, err, "update: ")
	})

	t.Run("Given a note present in the system, I can update its title", func(t *testing.T) {
		notesR, err := note.NewNotesRepo(fixtureNotes())
		assert.NoError(t, err)
		notesS := svc.NewNotesService(notesR)

		type testCase struct {
			noteID   uuid.UUID
			newTitle string
			want     note.Note
		}

		testCases := []testCase{
			{
				noteID:   uuid.UUID{1},
				newTitle: "New title",
				want: note.Note{
					NoteID:  uuid.UUID{1},
					Title:   "New title",
					Content: "robs 1st note content",
					UserID:  uuid.UUID{1},
				},
			},
		}

		for _, tc := range testCases {
			got, _ := notesS.Update(tc.noteID, tc.newTitle)
			assert.Equal(t, tc.want, got)

			got = notesR.GetNoteByID(tc.noteID)
			assert.Equal(t, tc.want, got)
		}
	})

	t.Run("I can get a note by its ID", func(t *testing.T) {
		notesR, err := note.NewNotesRepo(fixtureNotes())
		assert.NoError(t, err)
		type testCase struct {
			noteID uuid.UUID
			want   note.Note
		}

		testCases := []testCase{
			{
				noteID: uuid.UUID{1}, want: note.Note{
					NoteID:  uuid.UUID{1},
					Title:   "robs 1st note",
					Content: "robs 1st note content",
					UserID:  uuid.UUID{1},
				},
			},
		}

		for _, tc := range testCases {
			got := notesR.GetNoteByID(tc.noteID)
			assert.Equal(t, tc.want, got)
		}
	})

	t.Run("I can get all notes of a User by the userID", func(t *testing.T) {
		notesR, err := note.NewNotesRepo(fixtureNotes())
		assert.NoError(t, err)
		type testCase struct {
			userID uuid.UUID
			want   []note.Note
		}

		testCases := []testCase{
			{userID: uuid.UUID{1},
				want: []note.Note{
					{UserID: uuid.UUID{1}, Title: "robs 1st note", Content: "robs 1st note content"},
					{UserID: uuid.UUID{1}, Title: "robs 2nd note", Content: "robs 2nd note content"},
				},
			},
			{userID: uuid.UUID{2},
				want: []note.Note{
					{UserID: uuid.UUID{2}, Title: "annas 1st note", Content: "annas 1st note content"},
					{UserID: uuid.UUID{2}, Title: "annas 2nd note", Content: "annas 2nd note content"},
				},
			},
		}

		for _, tc := range testCases {
			got := notesR.GetNotesByUserID(tc.userID)
			assert.ElementsMatch(t, tc.want, got)
		}
	})

	t.Run("When initialising NewNotesRepo, every note needs to have a different noteID", func(t *testing.T) {
		notes := []note.Note{
			{
				NoteID:  uuid.UUID{},
				UserID:  uuid.UUID{2},
				Title:   "annas 1st note",
				Content: "annas 1st note content",
			},
			{
				NoteID:  uuid.UUID{},
				UserID:  uuid.UUID{2},
				Title:   "annas 2nd note",
				Content: "annas 2nd note content",
			},
		}
		_, err := note.NewNotesRepo(notes)
		assert.ErrorContains(t, err, "newNotesRepo: duplicate noteID")
	})
}

func fixtureNotes() []note.Note {
	return []note.Note{
		{
			NoteID:  uuid.UUID{1},
			UserID:  uuid.UUID{1},
			Title:   "robs 1st note",
			Content: "robs 1st note content",
		},
		{

			NoteID:  uuid.UUID{2},
			UserID:  uuid.UUID{1},
			Title:   "robs 2nd note",
			Content: "robs 2nd note content",
		},
		{
			NoteID:  uuid.UUID{3},
			UserID:  uuid.UUID{2},
			Title:   "annas 1st note",
			Content: "annas 1st note content",
		},
		{
			NoteID:  uuid.UUID{4},
			UserID:  uuid.UUID{2},
			Title:   "annas 2nd note",
			Content: "annas 2nd note content",
		},
	}
}
