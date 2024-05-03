package note

import (
	"context"
	"fmt"

	"github.com/Keisn1/note-taking-app/domain/core/user"
	"github.com/google/uuid"
)

type Service interface {
	Delete(noteID uuid.UUID) error
	Create(nN UpdateNote) (Note, error)
	Update(n Note, newN UpdateNote) (Note, error)
	QueryByID(ctx context.Context, noteID uuid.UUID) (Note, error)
	GetNotesByUserID(userID uuid.UUID) ([]Note, error)
}

type NotesService struct {
	repo    Repo
	userSvc user.Service
}

func NewNotesService(nR Repo, us user.Service) NotesService {
	return NotesService{repo: nR, userSvc: us}
}

func (ns NotesService) Delete(noteID uuid.UUID) error {
	err := ns.repo.Delete(noteID)
	if err != nil {
		return fmt.Errorf("delete: [%s]", noteID)
	}
	return nil
}

func (ns NotesService) Create(ctx context.Context, nN UpdateNote) (Note, error) {
	// MidAuthenticate authenticates user but could still submit
	// a note with a UserID different from its id
	if _, err := ns.userSvc.QueryByID(ctx, nN.UserID); err != nil {
		return Note{}, err
	}

	n := Note{
		NoteID:  uuid.New(),
		Title:   nN.Title,
		Content: nN.Content,
		UserID:  nN.UserID,
	}

	err := ns.repo.Create(n)
	if err != nil {
		return Note{}, err
	}
	return n, nil
}

func (ns NotesService) Update(n Note, newN UpdateNote) (Note, error) {
	if !newN.Title.IsEmpty() {
		n.Title = newN.Title
	}

	if !newN.Content.IsEmpty() {
		n.Content = newN.Content
	}

	err := ns.repo.Update(n)
	if err != nil {
		return Note{}, fmt.Errorf("update: %w", err)
	}
	return n, nil
}

func (nS NotesService) QueryByID(ctx context.Context, noteID uuid.UUID) (Note, error) {
	n, err := nS.repo.QueryByID(ctx, noteID)
	if err != nil {
		return Note{}, fmt.Errorf("getNoteByID: [%s]: %w", noteID, err)
	}
	return n, nil
}

func (nS NotesService) GetNotesByUserID(userID uuid.UUID) ([]Note, error) {
	notes, err := nS.repo.GetNotesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("getNoteByUserID: [%s]: %w", userID, err)
	}
	return notes, nil
}
