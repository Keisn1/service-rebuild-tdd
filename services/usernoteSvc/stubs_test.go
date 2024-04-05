package usernoteSvc_test

import (
	"fmt"

	"github.com/Keisn1/note-taking-app/domain/user"
	"github.com/Keisn1/note-taking-app/domain/usernote"
	"github.com/google/uuid"
)

type StubUserRepository struct {
	users map[uuid.UUID]user.User
}

func (sUR *StubUserRepository) GetUserByID(userID uuid.UUID) (user.User, error) {
	u, ok := sUR.users[userID]
	if !ok {
		return user.User{}, fmt.Errorf("GetUserByID: user not found [%s]", userID)
	}
	return u, nil
}

type StubUserNoteRepository struct {
	Usernotes     map[uuid.UUID]usernote.UserNote
	EditWasCalled bool
}

func (sUNR *StubUserNoteRepository) Create(userID uuid.UUID, title, content string) (usernote.UserNote, error) {
	u := usernote.NewUserNote(title, content, userID)
	sUNR.Usernotes[u.GetID()] = u
	return u, nil
}

func (sUNR *StubUserNoteRepository) GetNoteByID(noteID uuid.UUID) (usernote.UserNote, error) {
	n, ok := sUNR.Usernotes[noteID]
	if !ok {
		return usernote.UserNote{}, fmt.Errorf("Note note found")
	}
	return n, nil
}

func (sUNR *StubUserNoteRepository) GetNotesByUserID(userID uuid.UUID) ([]usernote.UserNote, error) {
	var ret []usernote.UserNote
	for _, n := range sUNR.Usernotes {
		if n.GetUserID() == userID {
			ret = append(ret, n)
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("No note found for userID[%s]", userID)
	}
	return ret, nil
}
