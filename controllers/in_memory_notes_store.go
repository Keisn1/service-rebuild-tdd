package controllers

type InMemoryNotesStore struct {
	notes map[int][]string
}

func NewInMemoryNotesStore(data map[int][]string) InMemoryNotesStore {
	return InMemoryNotesStore{notes: data}
}

func (i *InMemoryNotesStore) GetNotesByID(id int) []string {
	return i.notes[id]
}

func (i *InMemoryNotesStore) GetAllNotes() []string {
	var allNotes []string
	for _, notes := range i.notes {
		allNotes = append(allNotes, notes...)
	}
	return allNotes
}

func (i *InMemoryNotesStore) AddNote(userID int, note string) error {
	i.notes[userID] = append(i.notes[userID], note)
	return nil
}
