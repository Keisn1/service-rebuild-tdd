package controllers

import (
	"fmt"
	"github.com/Keisn1/note-taking-app/domain"
	"sync"
)

type InMemoryNotesStore struct {
	notes domain.Notes
	lock  sync.RWMutex
}

func NewInMemoryNotesStore() *InMemoryNotesStore {
	return &InMemoryNotesStore{
		notes: domain.Notes{},
		lock:  sync.RWMutex{},
	}
}

func (i *InMemoryNotesStore) Delete(userID, noteID int) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	// Find the index of the note with the given userID and noteID
	index := -1
	for idx, note := range i.notes {
		if note.UserID == userID && note.NoteID == noteID {
			index = idx
			break
		}
	}

	// If the note is not found, return an error
	if index == -1 {
		return fmt.Errorf("note with UserID %d and NoteID %d not found", userID, noteID)
	}

	// Delete the note by slicing it out of the slice
	i.notes = append(i.notes[:index], i.notes[index+1:]...)

	return nil
}

// Update edits the note with the given userID and noteID with the new content
func (i *InMemoryNotesStore) EditNote(userID, noteID int, newNote string) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	// Find the note with the given userID and noteID
	for idx, note := range i.notes {
		if note.UserID == userID && note.NoteID == noteID {
			// Update the note with the new content
			i.notes[idx].Note = newNote
			return nil
		}
	}

	// If the note is not found, return an error
	return fmt.Errorf("note with UserID %d and NoteID %d not found", userID, noteID)
}

func (i *InMemoryNotesStore) GetNoteByUserIDAndNoteID(userID, noteID int) (domain.Notes, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	fmt.Println(i.notes)
	var userNotes domain.Notes
	for _, note := range i.notes {
		if note.UserID == userID && note.NoteID == noteID {
			userNotes = append(userNotes, note)
		}
	}
	return userNotes, nil
}

func (i *InMemoryNotesStore) GetNotesByUserID(userID int) (domain.Notes, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	var userNotes domain.Notes
	for _, note := range i.notes {
		if note.UserID == userID {
			userNotes = append(userNotes, note)
		}
	}
	return userNotes, nil
}

func (i *InMemoryNotesStore) GetAllNotes() (domain.Notes, error) {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.notes, nil
}

func (i *InMemoryNotesStore) AddNote(userID int, note string) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	noteID := len(i.notes) + 1
	newNote := domain.Note{NoteID: noteID, UserID: userID, Note: note}
	i.notes = append(i.notes, newNote)
	return nil
}
