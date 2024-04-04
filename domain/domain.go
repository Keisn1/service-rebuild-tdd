package domain

import (
	"fmt"

	"github.com/Keisn1/note-taking-app/domain/usernote"
	"github.com/google/uuid"
)

type UserNoteRepository interface {
	GetNoteByID(noteID uuid.UUID) (usernote.UserNote, error)
	GetNotesByUserID(userID uuid.UUID) ([]usernote.UserNote, error)
	Create(userID uuid.UUID, title, content string) (usernote.UserNote, error)
}

type UserNoteService struct {
	usernotes UserNoteRepository
}

func NewUserNoteService(cfgs ...ServiceConfig) UserNoteService {
	s := UserNoteService{}
	for _, cfg := range cfgs {
		cfg(&s)
	}
	return s
}

type ServiceConfig func(*UserNoteService) error

func WithUserNoteRepository(u UserNoteRepository) ServiceConfig {
	return func(s *UserNoteService) error {
		s.usernotes = u
		return nil
	}
}

func (s UserNoteService) Create(uID uuid.UUID, title, content string) (usernote.UserNote, error) {
	u, err := s.usernotes.Create(uID, title, content)
	return u, err
}

func (s UserNoteService) QueryByID(nID uuid.UUID) (usernote.UserNote, error) {
	n, err := s.usernotes.GetNoteByID(nID)
	if err != nil {
		return usernote.UserNote{}, fmt.Errorf("querybyid: noteID[%s]: %w", nID, err)
	}
	return n, nil
}

func (s UserNoteService) QueryByUserID(uID uuid.UUID) ([]usernote.UserNote, error) {
	notes, err := s.usernotes.GetNotesByUserID(uID)
	if err != nil {
		return nil, fmt.Errorf("querybyuserid: userID[%s]: %w", uID, err)
	}
	return notes, err
}
