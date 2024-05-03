package note_test

import (
	"testing"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/Keisn1/note-taking-app/domain/core/note/repositories/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func fixtureNotes() []note.Note {
	return []note.Note{
		note.NewNote(uuid.UUID{1}, "robs 1st note", "robs 1st note content", uuid.UUID{1}),
		note.NewNote(uuid.UUID{2}, "robs 2nd note", "robs 2nd note content", uuid.UUID{1}),
		note.NewNote(uuid.UUID{3}, "annas 1st note", "annas 1st note content", uuid.UUID{2}),
		note.NewNote(uuid.UUID{4}, "annas 2nd note", "annas 2nd note content", uuid.UUID{2}),
	}
}

func Setup(t *testing.T, notes []note.Note) note.NotesService {
	t.Helper()
	repo, err := memory.NewRepo(notes)
	assert.NoError(t, err)

	userSvc := StubUserService{ids: make(map[uuid.UUID]struct{})}
	for _, n := range notes {
		userSvc.ids[n.UserID] = struct{}{}
	}

	return note.NewNotesService(repo, userSvc)
}
