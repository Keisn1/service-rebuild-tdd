package controllers

// import "sync"

// type InMemoryNotesStore struct {
// 	notes map[int]Note
// 	lock  sync.RWMutex
// }

// func NewInMemoryNotesStore() InMemoryNotesStore {
// 	return InMemoryNotesStore{notes: make(map[int]Note), lock: sync.RWMutex{}}
// }

// func (i *InMemoryNotesStore) Delete(id int) error {
// 	return nil
// }

// func (i *InMemoryNotesStore) EditNote(note Note) error {
// 	return nil
// }
// func (i *InMemoryNotesStore) GetNotesByUserID(userID int) (ret Notes) {
// 	i.lock.Lock()
// 	defer i.lock.Unlock()

// 	for _, n := range i.notes {
// 		if n.UserID == userID {
// 			ret = append(ret, n)
// 		}
// 	}
// 	return
// }

// func (i *InMemoryNotesStore) GetAllNotes() Notes {
// 	i.lock.Lock()
// 	defer i.lock.Unlock()
// 	var allNotes Notes
// 	for _, note := range i.notes {
// 		allNotes = append(allNotes, note)
// 	}
// 	return allNotes
// }

// func (i *InMemoryNotesStore) AddNote(note Note) error {
// 	i.lock.Lock()
// 	defer i.lock.Unlock()
// 	id := len(i.notes) + 1
// 	i.notes[id] = note
// 	return nil
// }
