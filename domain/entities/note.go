package entities

import "github.com/google/uuid"

type Note struct {
	ID      NoteID
	Title   Title
	Content Content
}

type NoteID uuid.UUID

type Title string

type Content string
