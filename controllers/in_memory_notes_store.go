package controllers

import "sync"

type InMemoryNotesStore struct {
	notes map[int][]string
	lock  sync.RWMutex
}

func NewInMemoryNotesStore(data map[int][]string) InMemoryNotesStore {
	return InMemoryNotesStore{notes: data, lock: sync.RWMutex{}}
}

func (i *InMemoryNotesStore) GetNotesByID(id int) []string {
	i.lock.Lock()
	defer i.lock.Unlock()

	return i.notes[id]
}

func (i *InMemoryNotesStore) GetAllNotes() []string {
	i.lock.Lock()
	defer i.lock.Unlock()

	var allNotes []string
	for _, notes := range i.notes {
		allNotes = append(allNotes, notes...)
	}
	return allNotes
}

func (i *InMemoryNotesStore) AddNote(userID int, note string) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.notes[userID] = append(i.notes[userID], note)
	return nil
}
