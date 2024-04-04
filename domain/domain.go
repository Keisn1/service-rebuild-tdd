package domain

import (
	"fmt"

	"github.com/Keisn1/note-taking-app/domain/usernote"
	"github.com/google/uuid"
)

type UserNoteRepository interface {
	GetNoteByID(noteID uuid.UUID) (usernote.UserNote, error)
	GetNotesByUserID(userID uuid.UUID) ([]usernote.UserNote, error)
}

type Service struct {
	usernotes UserNoteRepository
}

func NewService(cfgs ...ServiceConfig) Service {
	s := Service{}
	for _, cfg := range cfgs {
		cfg(&s)
	}
	return s
}

type ServiceConfig func(*Service) error

func WithUserNoteRepository(u UserNoteRepository) ServiceConfig {
	return func(s *Service) error {
		s.usernotes = u
		return nil
	}
}

func (s Service) QueryByID(nID uuid.UUID) (usernote.UserNote, error) {
	n, err := s.usernotes.GetNoteByID(nID)
	if err != nil {
		return usernote.UserNote{}, fmt.Errorf("querybyid: noteID[%s]: %w", nID, err)
	}
	return n, nil
}

func (s Service) QueryByUserID(uID uuid.UUID) ([]usernote.UserNote, error) {
	notes, err := s.usernotes.GetNotesByUserID(uID)
	if err != nil {
		return nil, fmt.Errorf("querybyuserid: userID[%s]: %w", uID, err)
	}
	return notes, err
}
