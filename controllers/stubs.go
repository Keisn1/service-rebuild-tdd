package controllers

import (
	"errors"
	"fmt"
)

type AddNoteCall struct {
	userID int
	note   string
}

type DeleteCall struct {
	userID int
	noteID int
}

type EditCall struct {
	userID int
	noteID int
	note   string
}

type StubNotesStore struct {
	notes Notes

	getNotesByUserIDCalls []int
	allNotesGotCalled     bool
	addNoteCalls          []AddNoteCall
	editNoteCalls         []EditCall
	deleteNoteCalls       []DeleteCall
}

type StubNotesStoreFailureGetAllNotes struct {
	StubNotesStore
}

func (snsF *StubNotesStoreFailureGetAllNotes) GetAllNotes() (Notes, error) {
	return nil, DBError
}

func NewStubNotesStore() *StubNotesStore {
	return &StubNotesStore{
		notes: Notes{
			{NoteID: 1, UserID: 1, Note: "Note 1 user 1"},
			{NoteID: 2, UserID: 1, Note: "Note 2 user 1"},
			{NoteID: 3, UserID: 2, Note: "Note 1 user 2"},
			{NoteID: 4, UserID: 2, Note: "Note 2 user 2"},
		},
	}
}

func (sns *StubNotesStore) Delete(userID int, noteID int) error {
	sns.deleteNoteCalls = append(sns.deleteNoteCalls, DeleteCall{userID: userID, noteID: noteID})
	if userID == 50 && noteID == 50 {
		return errors.New("Resource not found")
	}
	return nil
}

func (sns *StubNotesStore) AddNote(userID int, note string) error {
	call := AddNoteCall{userID: userID, note: note}
	sns.addNoteCalls = append(sns.addNoteCalls, call)
	if call.userID == 1 && call.note == "Note already present" {
		return errors.New("Resource already exists")
	}
	return nil
}

func (sns *StubNotesStore) EditNote(userID, noteID int, note string) error {
	sns.editNoteCalls = append(sns.editNoteCalls, EditCall{userID: userID, noteID: noteID, note: note})
	return nil
}

func (sns *StubNotesStore) GetAllNotes() (Notes, error) {
	sns.allNotesGotCalled = true
	return nil, nil
}

func (sns *StubNotesStore) GetNotesByUserID(userID int) (ret Notes, err error) {
	sns.getNotesByUserIDCalls = append(sns.getNotesByUserIDCalls, userID)
	if userID == -1 {
		err = DBError
	}
	for _, n := range sns.notes {
		if n.UserID == userID {
			ret = append(ret, n)
		}
	}
	return
}

type StubLogger struct {
	infofCalls []string
	errorfCall []string
}

func (sl *StubLogger) Infof(format string, a ...any) {
	sl.infofCalls = append(sl.infofCalls, fmt.Sprintf(format, a...))
}

func (sl *StubLogger) Errorf(format string, a ...any) {
	sl.errorfCall = append(sl.errorfCall, fmt.Sprintf(format, a...))
}

func (sl *StubLogger) Reset() {
	sl.infofCalls = []string{}
	sl.errorfCall = []string{}
}

func NewStubLogger() *StubLogger {
	return &StubLogger{}
}
