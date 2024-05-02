package notesgrp_test

import (
	"github.com/Keisn1/note-taking-app/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockNotesSvc struct {
	mock.Mock
}

type mockNotesStoreParams struct {
	method          string
	arguments       []any
	returnArguments []any
}

func (mNS *mockNotesSvc) Setup(p mockNotesStoreParams) {
	mNS.Reset()
	mNS.On(p.method, p.arguments...).Return(p.returnArguments...)
}

func (mNS *mockNotesSvc) Reset() {
	mNS.Calls = []mock.Call{}
	mNS.ExpectedCalls = []*mock.Call{}
}

func (mNS *mockNotesSvc) GetAllNotes() (domain.Notes, error) {
	args := mNS.Called()
	return args.Get(0).(domain.Notes), args.Error(1)
}

func (mNS *mockNotesSvc) GetNoteByUserIDAndNoteID(userID uuid.UUID, noteID int) (domain.Notes, error) {
	args := mNS.Called(userID, noteID)
	return args.Get(0).(domain.Notes), args.Error(1)
}

func (mNS *mockNotesSvc) GetNotesByUserID(userID uuid.UUID) (domain.Notes, error) {
	args := mNS.Called(userID)
	return args.Get(0).(domain.Notes), args.Error(1)
}

func (mNS *mockNotesSvc) AddNote(userID uuid.UUID, note string) error {
	args := mNS.Called(userID, note)
	return args.Error(0)
}

func (mNS *mockNotesSvc) EditNote(userID uuid.UUID, noteID int, note string) error {
	return nil
}

func (mNS *mockNotesSvc) Delete(userID uuid.UUID, noteID int) error {
	args := mNS.Called(userID, noteID)
	return args.Error(0)
}
