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

func GetNotesByUserID(userID uuid.UUID) []Note {
	u := uuid.UUID{1}
	if userID == u {
		return []Note{
			{
				UserID:  uuid.UUID{1},
				Title:   "robs 1st note",
				Content: "robs 1st note content",
			},
			{
				UserID:  uuid.UUID{1},
				Title:   "robs 2nd note",
				Content: "robs 2nd note content",
			},
		}
	} else {
		return []Note{
			{
				UserID:  uuid.UUID{2},
				Title:   "annas 1st note",
				Content: "annas 1st note content",
			},
			{
				UserID:  uuid.UUID{2},
				Title:   "annas 2nd note",
				Content: "annas 2nd note content",
			},
		}

	}
}
