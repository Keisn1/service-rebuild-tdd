// Package usernoteSvc wraps the data/store layer
// handles Crud operations on aggregate usernote
// make changes persistent by calling data/store layer
package usernoteSvc

import "github.com/google/uuid"

type Note struct {
	NoteID  uuid.UUID
	Title   string
	Content string
	UserID  uuid.UUID
}

type notesRepo struct {
	notes []Note
}

func NewNotesRepo(notes []Note) notesRepo {
	var nR notesRepo
	nR.notes = notes
	return nR
}

func (nR notesRepo) GetNoteByID(noteID uuid.UUID) Note {
	for _, n := range nR.notes {
		if n.NoteID == noteID {
			return n
		}
	}
	return Note{}
}
func (nR notesRepo) GetNotesByUserID(userID uuid.UUID) []Note {
	var ret []Note
	for _, n := range nR.notes {
		if n.UserID == userID {
			n.NoteID = uuid.UUID{0}
			ret = append(ret, n)
		}
	}
	return ret
}
