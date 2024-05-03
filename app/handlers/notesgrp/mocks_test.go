package notesgrp_test

import (
	"github.com/Keisn1/note-taking-app/domain/core/note"
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

func (mNS *mockNotesSvc) GetNotesByUserID(userID uuid.UUID) ([]note.Note, error) {
	args := mNS.Called(userID)
	return args.Get(0).([]note.Note), args.Error(1)
}

func (mNS *mockNotesSvc) QueryByID(noteID uuid.UUID) (note.Note, error) {
	args := mNS.Called(noteID)
	return args.Get(0).(note.Note), args.Error(1)
}

func (mNS *mockNotesSvc) Create(newN note.UpdateNote) (note.Note, error) {
	args := mNS.Called(newN)
	return args.Get(0).(note.Note), args.Error(1)
}

func (mNS *mockNotesSvc) Update(n note.Note, un note.UpdateNote) (note.Note, error) {
	return note.Note{}, nil
}

func (mNS *mockNotesSvc) Delete(noteID uuid.UUID) error {
	args := mNS.Called(noteID)
	return args.Error(0)
}
