package entities

import "github.com/google/uuid"

type Person struct {
	ID    uuid.UUID
	Name  Username
	Email string
}

type Username string

type Email string
