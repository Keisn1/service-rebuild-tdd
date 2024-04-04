package entities

import "github.com/google/uuid"

type Person struct {
	ID    uuid.UUID
	Name  Username
	Email Email
}

type Username string

type Email string
