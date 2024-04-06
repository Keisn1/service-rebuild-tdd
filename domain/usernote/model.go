package usernote

import (
	ents "github.com/Keisn1/note-taking-app/domain/entities"
	"github.com/google/uuid"
)

type UserNote struct {
	note   *ents.Note
	userID ents.UserID
}

func NewUserNote(title, content string, userID uuid.UUID) UserNote {
	return UserNote{
		note: &ents.Note{
			ID:      ents.NoteID(uuid.New()),
			Title:   ents.Title(title),
			Content: ents.Content(content),
		},
		userID: ents.UserID(userID),
	}
}

func (u UserNote) GetTitle() ents.Title {
	return u.note.Title
}

func (u UserNote) SetTitle(title string) {
	u.note.Title = ents.Title(title)
}

func (u UserNote) GetContent() ents.Content {
	return u.note.Content
}

func (u UserNote) SetContent(content string) {
	u.note.Content = ents.Content(content)
}

func (u UserNote) GetID() ents.NoteID {
	return u.note.ID
}

func (u UserNote) SetID(noteID uuid.UUID) {
	u.note.ID = ents.NoteID(noteID)
}

func (u UserNote) GetUserID() ents.UserID {
	return u.userID
}
