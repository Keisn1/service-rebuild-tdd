package controllers_test

import (
	"github.com/Keisn1/note-taking-app/domain"
	"github.com/stretchr/testify/mock"
)

type mockNotesStore struct {
	mock.Mock
}

func (mNS *mockNotesStore) GetAllNotes() (domain.Notes, error) {
	args := mNS.Called()
	return args.Get(0).(domain.Notes), args.Error(1)
}

func (mNS *mockNotesStore) GetNoteByUserIDAndNoteID(userID, noteID int) (domain.Notes, error) {
	return nil, nil
}

func (mNS *mockNotesStore) GetNotesByUserID(userID int) (domain.Notes, error) {
	return nil, nil
}

func (mNS *mockNotesStore) AddNote(userID int, note string) error {
	return nil
}

func (mNS *mockNotesStore) EditNote(userID, noteID int, note string) error {
	return nil
}

func (mNS *mockNotesStore) Delete(userID, noteID int) error {
	return nil
}

type mockLogger struct {
	mock.Mock
}

func (ml *mockLogger) Infof(format string, args ...any) {
	if len(args) == 0 {
		ml.Called(format)
	} else {
		ml.Called(format, args)
	}
}

func (ml *mockLogger) Errorf(format string, a ...any) {
	if len(a) == 0 {
		ml.Called(format)
	} else {
		ml.Called(format, a)
	}
}
