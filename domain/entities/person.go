package entities

import (
	"net/mail"

	"github.com/google/uuid"
)

type Person struct {
	ID    uuid.UUID
	Name  Username
	Email mail.Address
}

type Username string
