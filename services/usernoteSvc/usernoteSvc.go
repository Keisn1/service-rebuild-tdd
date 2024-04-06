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

type notesRepo struct {
	notes []Note
}

func NewNotesRepo(notes []Note) notesRepo {
	var nR notesRepo
	nR.notes = notes
	return nR
}

func (nR notesRepo) GetNotesByUserID(userID uuid.UUID) []Note {
	var ret []Note
	for _, n := range nR.notes {
		if n.UserID == userID {
			ret = append(ret, n)
		}
	}
	return ret
}
