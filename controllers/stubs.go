package controllers

import (
	"errors"
)

type StubNotesStore struct {
	notes map[int]Note

	getNotesByUserIDCalls []int
	allNotesGotCalled     bool
	addNoteCalls          Notes
	editNoteCalls         Notes
	deleteNoteCalls       []int
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
	sns.deleteNoteCalls = append(sns.deleteNoteCalls, id)
	if id != 1 {
		return errors.New("Resource not found")
	}
	return nil
}

func (sns *StubNotesStore) AddNote(note Note) error {
	sns.addNoteCalls = append(sns.addNoteCalls, note)
	if note.UserID == 1 && note.Note == "Note 1 user 1" {
		return errors.New("Resource already exists")
	}
	return nil
}

func (sns *StubNotesStore) EditNote(note Note) error {
	sns.editNoteCalls = append(sns.editNoteCalls, note)
	return nil
}

func (sns *StubNotesStore) GetAllNotes() Notes {
	sns.allNotesGotCalled = true
	return nil
}

func (sns *StubNotesStore) GetNotesByUserID(userID int) (ret Notes) {
	for _, n := range sns.notes {
		if n.UserID == userID {
			ret = append(ret, n)
		}
	}
	return
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
