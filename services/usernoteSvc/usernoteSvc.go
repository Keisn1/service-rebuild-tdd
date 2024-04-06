// Package usernoteSvc wraps the data/store layer
// handles Crud operations on aggregate usernote
// make changes persistent by calling data/store layer
package usernoteSvc

import "github.com/google/uuid"

type Note struct {
	UserID  uuid.UUID
	Title   string
	Content string
}

func GetNoteByUserID(userID uuid.UUID) Note {
	return Note{
		UserID:  userID,
		Title:   "robs note",
		Content: "robs note content",
	}
}
