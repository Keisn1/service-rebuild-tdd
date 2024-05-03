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
		{ID: uuid.UUID{1}, Title: note.NewTitle("robs 1st note"), Content: note.NewContent("robs 1st note content"), UserID: uuid.UUID{1}},
		{ID: uuid.UUID{2}, Title: note.NewTitle("robs 2nd note"), Content: note.NewContent("robs 2nd note content"), UserID: uuid.UUID{1}},
		{ID: uuid.UUID{3}, Title: note.NewTitle("annas 1st note"), Content: note.NewContent("annas 1st note content"), UserID: uuid.UUID{2}},
		{ID: uuid.UUID{4}, Title: note.NewTitle("annas 2nd note"), Content: note.NewContent("annas 2nd note content"), UserID: uuid.UUID{2}},
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
