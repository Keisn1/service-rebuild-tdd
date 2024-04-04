package entities

import "github.com/google/uuid"

type Note struct {
	ID      uuid.UUID
	Title   Title
	Content Content
}

type Title string

type Content string

func NewNote(id uuid.UUID, title, content string) *Note {
	return &Note{
		ID:      id,
		Title:   Title(title),
		Content: Content(content),
	}
}
