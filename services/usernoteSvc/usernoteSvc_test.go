package usernoteSvc_test

import (
	"testing"

	svc "github.com/Keisn1/note-taking-app/services/usernoteSvc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNotes(t *testing.T) {
	t.Run("When initialising NewNotesRepo, every note needs to have a different noteID", func(t *testing.T) {
		notes := []svc.Note{
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
		_, err := svc.NewNotesRepo(notes)
		assert.ErrorContains(t, err, "newNotesRepo: duplicate noteID")
	})

	t.Run("I can get a note by its ID", func(t *testing.T) {
		unRepo, err := svc.NewNotesRepo(fixtureNotes())
		assert.NoError(t, err)
		type testCase struct {
			noteID uuid.UUID
			want   svc.Note
		}

		testCases := []testCase{
			{
				noteID: uuid.UUID{1}, want: svc.Note{
					NoteID:  uuid.UUID{1},
					Title:   "robs 1st note",
					Content: "robs 1st note content",
					UserID:  uuid.UUID{1},
				},
			},
		}

		for _, tc := range testCases {
			got := unRepo.GetNoteByID(tc.noteID)
			assert.Equal(t, tc.want, got)
		}
	})

	t.Run("I can get all notes of a User by the userID", func(t *testing.T) {
		unRepo, err := svc.NewNotesRepo(fixtureNotes())
		assert.NoError(t, err)
		type testCase struct {
			userID uuid.UUID
			want   []svc.Note
		}

		testCases := []testCase{
			{userID: uuid.UUID{1},
				want: []svc.Note{
					{UserID: uuid.UUID{1}, Title: "robs 1st note", Content: "robs 1st note content"},
					{UserID: uuid.UUID{1}, Title: "robs 2nd note", Content: "robs 2nd note content"},
				},
			},
			{userID: uuid.UUID{2},
				want: []svc.Note{
					{UserID: uuid.UUID{2}, Title: "annas 1st note", Content: "annas 1st note content"},
					{UserID: uuid.UUID{2}, Title: "annas 2nd note", Content: "annas 2nd note content"},
				},
			},
		}

		for _, tc := range testCases {
			got := unRepo.GetNotesByUserID(tc.userID)
			assert.Equal(t, tc.want, got)
		}
	})
}

func fixtureNotes() []svc.Note {
	return []svc.Note{
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
