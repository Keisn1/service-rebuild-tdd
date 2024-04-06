package entities

import (
	"net/mail"

	"github.com/google/uuid"
)

type Person struct {
	ID    UserID
	Name  Username
	Email mail.Address
}

type UserID uuid.UUID

type Username string
