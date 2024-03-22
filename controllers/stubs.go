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
	Notes Notes

	getNotesByUserIDCalls         []int
	getNoteByUserIDAndNoteIDCalls [][2]int
	getAllNotesGotCalled          bool
	addNoteCalls                  []AddNoteCall
	editNoteCalls                 []EditCall
	deleteNoteCalls               []DeleteCall
}

type StubNotesStoreFailureGetAllNotes struct {
	StubNotesStore
}

func (snsF *StubNotesStoreFailureGetAllNotes) GetAllNotes() (Notes, error) {
	snsF.getAllNotesGotCalled = true
	return nil, ErrDB
}

func NewStubNotesStore() *StubNotesStore {
	return &StubNotesStore{
		Notes: Notes{
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
	sns.getAllNotesGotCalled = true
	return sns.Notes, nil
}

func (sns *StubNotesStore) GetNoteByUserIDAndNoteID(userID, noteID int) (Notes, error) {
	sns.getNoteByUserIDAndNoteIDCalls = append(sns.getNoteByUserIDAndNoteIDCalls, [2]int{userID, noteID})
	var userNotes Notes
	if userID == -1 {
		err := ErrDB
		return nil, err
	}
	for _, note := range sns.Notes {
		if note.UserID == userID && note.NoteID == noteID {
			userNotes = append(userNotes, note)
		}
	}
	return userNotes, nil
}

func (sns *StubNotesStore) GetNotesByUserID(userID int) (ret Notes, err error) {
	sns.getNotesByUserIDCalls = append(sns.getNotesByUserIDCalls, userID)
	if userID == -1 {
		err = ErrDB
		return nil, err
	}
	for _, n := range sns.Notes {
		if n.UserID == userID {
			ret = append(ret, n)
		}
	}
	return ret, err
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
