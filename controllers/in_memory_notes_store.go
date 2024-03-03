package controllers

import "sync"

type InMemoryNotesStore struct {
	notes Notes
	lock  sync.RWMutex
}

func NewInMemoryNotesStore() InMemoryNotesStore {
	return InMemoryNotesStore{notes: Notes{}, lock: sync.RWMutex{}}
}

func (i *InMemoryNotesStore) Delete(id int) error {
	return nil
}

func (i *InMemoryNotesStore) EditNote(note Note) error {
	return nil
}
func (i *InMemoryNotesStore) GetNotesByUserID(userID int) (ret Notes) {
	i.lock.Lock()
	defer i.lock.Unlock()

	for _, n := range i.notes {
		if n.UserID == userID {
			ret = append(ret, n)
		}
	}
	return
}

func (i *InMemoryNotesStore) GetAllNotes() Notes {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.notes
}

func (i *InMemoryNotesStore) AddNote(note Note) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.notes = append(i.notes, note)
	return nil
}
