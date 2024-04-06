package usernoteSvc_test

import (
	"testing"

	svc "github.com/Keisn1/note-taking-app/services/usernoteSvc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNotes(t *testing.T) {

	t.Run("I can get a note of a its ID", func(t *testing.T) {
		unRepo := svc.NewNotesRepo(fixtureNotes())
		noteID := uuid.UUID{1}
		got := unRepo.GetNoteByID(noteID)
		want := svc.Note{
			NoteID:  uuid.UUID{1},
			Title:   "robs 1st note",
			Content: "robs 1st content",
			UserID:  uuid.UUID{1},
		}
		assert.Equal(t, want, got)
	})

	t.Run("I can get all notes of a User by the userID", func(t *testing.T) {
		unRepo := svc.NewNotesRepo(fixtureNotes())
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
			UserID:  uuid.UUID{1},
			Title:   "robs 1st note",
			Content: "robs 1st note content",
		},
		{
			UserID:  uuid.UUID{1},
			Title:   "robs 2nd note",
			Content: "robs 2nd note content",
		},
		{
			UserID:  uuid.UUID{2},
			Title:   "annas 1st note",
			Content: "annas 1st note content",
		},
		{
			UserID:  uuid.UUID{2},
			Title:   "annas 2nd note",
			Content: "annas 2nd note content",
		},
	}
}
