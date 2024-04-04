package domain

import (
	"github.com/Keisn1/note-taking-app/domain/usernote"
	"github.com/google/uuid"
)

type UserNoteRepository interface {
	GetNoteByID(noteID uuid.UUID) (usernote.UserNote, error)
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

func (s Service) GetNoteByID(nID uuid.UUID) (usernote.UserNote, error) {
	n, err := s.usernotes.GetNoteByID(nID)
	if err != nil {
		return usernote.UserNote{}, err
	}
	return n, nil
}
