package entities

import "github.com/google/uuid"

type Note struct {
	ID      uuid.UUID
	Title   Title
	Content Content
}

type Title string

type Content string
