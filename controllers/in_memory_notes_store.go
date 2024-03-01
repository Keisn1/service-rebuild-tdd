package controllers

import "sync"

type InMemoryNotesStore struct {
	notes map[int]Notes
	lock  sync.RWMutex
}

func NewInMemoryNotesStore(data map[int]Notes) InMemoryNotesStore {
	return InMemoryNotesStore{notes: data, lock: sync.RWMutex{}}
}

func (i *InMemoryNotesStore) GetNotesByID(id int) Notes {
	i.lock.Lock()
	defer i.lock.Unlock()

	return i.notes[id]
}

func (i *InMemoryNotesStore) GetAllNotes() map[int]Notes {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.notes
}

func (i *InMemoryNotesStore) AddNote(userID int, note string) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.notes[userID] = append(i.notes[userID], note)
	return nil
}
