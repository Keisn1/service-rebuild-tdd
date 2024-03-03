package controllers

import (
	"errors"
	"fmt"
)

type StubNotesStore struct {
	notes         map[int]Note
	addNoteCalls  Notes
	editNoteCalls Notes
}

func NewStubNotesStore() *StubNotesStore {
	return &StubNotesStore{
		notes: map[int]Note{
			1: {UserID: 1, Note: "Note 1 user 1"},
			2: {UserID: 1, Note: "Note 2 user 1"},
			3: {UserID: 2, Note: "Note 1 user 2"},
			4: {UserID: 2, Note: "Note 2 user 2"},
		},
	}
}

func (sns *StubNotesStore) Delete(id int) error {
	if _, ok := sns.notes[id]; ok {
		delete(sns.notes, id)
		return nil
	}
	return errors.New(fmt.Sprintf("Note with id %v not found", id))
}

func (sns *StubNotesStore) AddNote(note Note) error {
	sns.addNoteCalls = append(sns.addNoteCalls, note)
	return nil
}

func (sns *StubNotesStore) EditNote(note Note) error {
	sns.editNoteCalls = append(sns.editNoteCalls, note)
	return nil
}

func (sns *StubNotesStore) GetAllNotes() Notes {
	var allNotes Notes
	for _, note := range sns.notes {
		allNotes = append(allNotes, note)
	}
	return allNotes
}

func (sns *StubNotesStore) GetNotesByUserID(userID int) (ret Notes) {
	for _, n := range sns.notes {
		if n.UserID == userID {
			ret = append(ret, n)
		}
	}
	return
}

type StubNotesStoreAddNoteErrors struct {
	StubNotesStore
}

func (sns *StubNotesStoreAddNoteErrors) AddNote(note Note) error {
	return errors.New("Error stub")
}

type fmtCallf struct {
	format string
	a      []any
}

type StubLogger struct {
	infofCalls []fmtCallf
	errorfCall []fmtCallf
}

func (sl *StubLogger) Infof(format string, a ...any) {
	sl.infofCalls = append(sl.infofCalls, fmtCallf{format: format, a: a})
}

func (sl *StubLogger) Errorf(format string, a ...any) {
	sl.errorfCall = append(sl.errorfCall, fmtCallf{format: format, a: a})
}

func (sl *StubLogger) Reset() {
	sl.infofCalls = []fmtCallf{}
	sl.errorfCall = []fmtCallf{}
}

func NewStubLogger() *StubLogger {
	return &StubLogger{}
}
