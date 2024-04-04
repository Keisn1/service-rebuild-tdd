package entities

import "github.com/google/uuid"

type Note struct {
	ID      uuid.UUID
	Title   Title
	Content Text
}

type Text string

type Title string
