package controllers_test

import (
	"github.com/Keisn1/note-taking-app/domain"
	"github.com/stretchr/testify/mock"
)

type mockNotesStore struct {
	mock.Mock
}

type mockNotesStoreParams struct {
	method          string
	arguments       []any
	returnArguments []any
}

func (mNS *mockNotesStore) Setup(p mockNotesStoreParams) {
	mNS.Reset()
	mNS.On(p.method, p.arguments...).Return(p.returnArguments...)
}

func (mNS *mockNotesStore) Reset() {
	mNS.Calls = []mock.Call{}
	mNS.ExpectedCalls = []*mock.Call{}
}

func (mNS *mockNotesStore) GetAllNotes() (domain.Notes, error) {
	args := mNS.Called()
	return args.Get(0).(domain.Notes), args.Error(1)
}

func (mNS *mockNotesStore) GetNoteByUserIDAndNoteID(userID, noteID int) (domain.Notes, error) {
	args := mNS.Called(userID, noteID)
	return args.Get(0).(domain.Notes), args.Error(1)
}

func (mNS *mockNotesStore) GetNotesByUserID(userID int) (domain.Notes, error) {
	args := mNS.Called(userID)
	return args.Get(0).(domain.Notes), args.Error(1)
}

func (mNS *mockNotesStore) AddNote(userID int, note string) error {
	args := mNS.Called(userID, note)
	return args.Error(0)
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

type mockLoggingParams struct {
	method    string
	arguments []any
}

func (ml *mockLogger) Setup(lParams mockLoggingParams) {
	ml.Reset()
	ml.On(lParams.method, lParams.arguments...).Return(nil)
}

func (ml *mockLogger) Reset() {
	ml.ExpectedCalls = []*mock.Call{}
	ml.Calls = []mock.Call{}
}

func (ml *mockLogger) Infof(format string, args ...any) {
	var arguments mock.Arguments
	arguments = append(arguments, format)
	arguments = append(arguments, args...)
	ml.Called(arguments...)
}

func (ml *mockLogger) Errorf(format string, a ...any) {
	var arguments mock.Arguments
	arguments = append(arguments, format)
	arguments = append(arguments, a...)
	ml.Called(arguments...)
}
