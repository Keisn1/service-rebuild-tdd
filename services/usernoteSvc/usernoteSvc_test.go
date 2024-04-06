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

		err = notesS.Update(note.Note{}, note.Note{})
		assert.ErrorContains(t, err, "update: ")
	})

	t.Run("Given a note present in the system, I can update its title and its content", func(t *testing.T) {
		notesR, err := note.NewNotesRepo(fixtureNotes())
		assert.NoError(t, err)
		notesS := svc.NewNotesService(notesR)

		type testCase struct {
			name       string
			currNote   note.Note
			updateNote note.Note
			want       note.Note
		}

		testCases := []testCase{
			{
				name:       "New Title, new 0 character content, update of both (title and content) expected",
				currNote:   note.NewNote(uuid.UUID{1}, note.NewTitle("robs 1st note"), note.NewContent("robs 1st note content"), uuid.UUID{1}),
				updateNote: note.NewNote(uuid.UUID{1}, note.NewTitle("New title"), note.NewContent(""), uuid.UUID{1}),
				want:       note.NewNote(uuid.UUID{1}, note.NewTitle("New title"), note.NewContent(""), uuid.UUID{1}),
			},
			{
				name:       "New Title, empty content, therefore no update",
				currNote:   note.NewNote(uuid.UUID{2}, note.NewTitle("robs 2nd note"), note.NewContent("robs 2nd note content"), uuid.UUID{1}),
				updateNote: note.NewNote(uuid.UUID{2}, note.NewTitle("New title"), note.Content{}, uuid.UUID{1}),
				want:       note.NewNote(uuid.UUID{2}, note.NewTitle("New title"), note.NewContent("robs 2nd note content"), uuid.UUID{1}),
			},
			{
				name:       "New 0 character title, new content, update of both (title and content) expected",
				currNote:   note.NewNote(uuid.UUID{3}, note.NewTitle("annas 1st note"), note.NewContent("annas 1st note content"), uuid.UUID{2}),
				updateNote: note.NewNote(uuid.UUID{3}, note.NewTitle(""), note.NewContent("New content"), uuid.UUID{2}),
				want:       note.NewNote(uuid.UUID{3}, note.NewTitle(""), note.NewContent("New content"), uuid.UUID{2}),
			},
			{
				name:       "empty title, empty content, therefore no update",
				currNote:   note.NewNote(uuid.UUID{4}, note.NewTitle("annas 2nd note"), note.NewContent("annas 2nd note content"), uuid.UUID{2}),
				updateNote: note.NewNote(uuid.UUID{4}, note.Title{}, note.Content{}, uuid.UUID{2}),
				want:       note.NewNote(uuid.UUID{4}, note.NewTitle("annas 2nd note"), note.NewContent("annas 2nd note content"), uuid.UUID{2}),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := notesS.Update(tc.currNote, tc.updateNote)
				assert.NoError(t, err)

				got := notesS.GetNoteByID(tc.currNote.GetID())
				assert.Equal(t, tc.want, got)
			})
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
			{noteID: uuid.UUID{1}, want: note.NewNote(uuid.UUID{1}, note.NewTitle("robs 1st note"), note.NewContent("robs 1st note content"), uuid.UUID{1})},
			{noteID: uuid.UUID{3}, want: note.NewNote(uuid.UUID{3}, note.NewTitle("annas 1st note"), note.NewContent("annas 1st note content"), uuid.UUID{2})},
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
			{
				userID: uuid.UUID{1},
				want: []note.Note{
					note.NewNote(uuid.UUID{}, note.NewTitle("robs 1st note"), note.NewContent("robs 1st note content"), uuid.UUID{1}),
					note.NewNote(uuid.UUID{}, note.NewTitle("robs 2nd note"), note.NewContent("robs 2nd note content"), uuid.UUID{1}),
				},
			},
			{
				userID: uuid.UUID{2},
				want: []note.Note{
					note.NewNote(uuid.UUID{}, note.NewTitle("annas 1st note"), note.NewContent("annas 1st note content"), uuid.UUID{2}),
					note.NewNote(uuid.UUID{}, note.NewTitle("annas 2nd note"), note.NewContent("annas 2nd note content"), uuid.UUID{2}),
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
			note.NewNote(uuid.UUID{}, note.NewTitle("annas 1st note"), note.NewContent("annas 1st note content"), uuid.UUID{2}),
			note.NewNote(uuid.UUID{}, note.NewTitle("annas 2nd note"), note.NewContent("annas 2nd note content"), uuid.UUID{2}),
		}
		_, err := note.NewNotesRepo(notes)
		assert.ErrorContains(t, err, "newNotesRepo: duplicate noteID")
	})
}

func fixtureNotes() []note.Note {
	return []note.Note{
		note.NewNote(uuid.UUID{1}, note.NewTitle("robs 1st note"), note.NewContent("robs 1st note content"), uuid.UUID{1}),
		note.NewNote(uuid.UUID{2}, note.NewTitle("robs 2nd note"), note.NewContent("robs 2nd note content"), uuid.UUID{1}),
		note.NewNote(uuid.UUID{3}, note.NewTitle("annas 1st note"), note.NewContent("annas 1st note content"), uuid.UUID{2}),
		note.NewNote(uuid.UUID{4}, note.NewTitle("annas 2nd note"), note.NewContent("annas 2nd note content"), uuid.UUID{2}),
	}
}
