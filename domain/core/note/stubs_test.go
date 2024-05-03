package note_test

import (
	"context"
	"errors"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/Keisn1/note-taking-app/domain/core/user"
	"github.com/google/uuid"
)

type ErrorNoteRepo struct {
	notes map[uuid.UUID]note.Note
}

func (nR ErrorNoteRepo) Create(n note.Note) error      { return errors.New("error in noteRepo") }
func (nR ErrorNoteRepo) Delete(noteID uuid.UUID) error { return nil }
func (nR ErrorNoteRepo) Update(note note.Note) error   { return nil }
func (nR ErrorNoteRepo) QueryByID(ctx context.Context, noteID uuid.UUID) (note.Note, error) {
	return note.Note{}, nil
}
func (nR ErrorNoteRepo) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) { return nil, nil }

type StubUserService struct {
	ids map[uuid.UUID]struct{}
}

func (sus StubUserService) QueryByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	if _, ok := sus.ids[userID]; !ok {
		return user.User{}, errors.New("User not found")
	}
	return user.User{ID: userID}, nil
}

func (sus StubUserService) Create(ctx context.Context, uu user.UpdateUser) (user.User, error) {
	return user.User{}, nil
}
func (sus StubUserService) Update(ctx context.Context, u user.User, uu user.UpdateUser) (user.User, error) {
	return user.User{}, nil
}
