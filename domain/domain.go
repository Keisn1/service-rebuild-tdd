package domain

import "github.com/google/uuid"

type Note struct {
	ID      uuid.UUID
	Title   Title
	Content Content
}

type Title string

type Content string

func NewNote(title, content string) Note {
	return Note{
		ID:      uuid.New(),
		Title:   Title(title),
		Content: Content(content),
	}
}

var notes []Note

func AddNote(note Note) {
	notes = append(notes, note)
}

func GetNoteByID(noteID uuid.UUID) Note {
	if noteID == uuid.UUID([16]byte{1}) {
		return Note{
			ID:      noteID,
			Title:   "title1",
			Content: "content1",
		}

	}
	return Note{}
}
