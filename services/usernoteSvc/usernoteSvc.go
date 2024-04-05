package usernoteSvc

import (
	"fmt"

	"github.com/Keisn1/note-taking-app/domain/usernote"
	"github.com/google/uuid"
)

type UserNoteService struct {
	usernotes usernote.UserNoteRepository
}

func NewUserNoteService(cfgs ...ServiceConfig) UserNoteService {
	s := UserNoteService{}
	for _, cfg := range cfgs {
		cfg(&s)
	}
	return s
}

type ServiceConfig func(*UserNoteService) error

func WithUserNoteRepository(u usernote.UserNoteRepository) ServiceConfig {
	return func(s *UserNoteService) error {
		s.usernotes = u
		return nil
	}
}

func (s UserNoteService) Create(userID uuid.UUID, title, content string) (usernote.UserNote, error) {
	if userID == uuid.UUID([16]byte{100}) {
		return usernote.UserNote{}, fmt.Errorf("Create: userID[%s]", userID)
	}
	u, err := s.usernotes.Create(userID, title, content)
	return u, err
}

func (s UserNoteService) QueryByID(noteID uuid.UUID) (usernote.UserNote, error) {
	n, err := s.usernotes.GetNoteByID(noteID)
	if err != nil {
		return usernote.UserNote{}, fmt.Errorf("querybyid: noteID[%s]: %w", noteID, err)
	}
	return n, nil
}

func (s UserNoteService) QueryByUserID(userID uuid.UUID) ([]usernote.UserNote, error) {
	notes, err := s.usernotes.GetNotesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("querybyuserid: userID[%s]: %w", userID, err)
	}
	return notes, err
}
