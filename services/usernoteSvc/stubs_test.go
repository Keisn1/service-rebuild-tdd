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

// func (sUR *StubUserRepository) GetUserByID(userID uuid.UUID) (user.User, error) {
// 	u, ok := sUR.users[userID]
// 	if !ok {
// 		return user.User{}, fmt.Errorf("GetUserByID: user not found [%s]", userID)
// 	}
// 	return u, nil
// }

type UserNoteDB struct {
	ID      uuid.UUID
	Title   string
	Content string
	UserID  uuid.UUID
}

type StubUserNoteRepository struct {
	Usernotes map[uuid.UUID]UserNoteDB
	// EditWasCalled bool
}

func toUserNoteDB(un usernote.UserNote) UserNoteDB {
	return UserNoteDB{
		ID:      uuid.UUID(un.GetID()),
		Title:   string(un.GetTitle()),
		Content: string(un.GetContent()),
		UserID:  uuid.UUID(un.GetUserID()),
	}
}

func toUserNote(unDB UserNoteDB) usernote.UserNote {
	un := usernote.NewUserNote(unDB.Title, unDB.Content, unDB.UserID)
	un.SetID(unDB.ID)
	return un
}

// func (sUNR *StubUserNoteRepository) Create(un usernote.UserNote) error {
// 	unDB := toUserNoteDB(un)
// 	sUNR.Usernotes[unDB.ID] = unDB
// 	return nil
// }

func (sUNR *StubUserNoteRepository) GetNoteByID(noteID uuid.UUID) (usernote.UserNote, error) {
	return usernote.UserNote{}, nil
	// noteDB := sUNR.Usernotes[noteID]
	// n, ok := sUNR.Usernotes[noteID]
	// if !ok {
	// 	return usernote.UserNote{}, fmt.Errorf("Note note found")
	// }

	// return n, nil
}

// func (sUNR *StubUserNoteRepository) GetNotesByUserID(userID uuid.UUID) ([]usernote.UserNote, error) {
// 	var ret []usernote.UserNote
// 	for _, n := range sUNR.Usernotes {
// 		if n.GetUserID() == userID {
// 			ret = append(ret, n)
// 		}
// 	}
// 	if len(ret) == 0 {
// 		return nil, fmt.Errorf("No note found for userID[%s]", userID)
// 	}
// 	return ret, nil
// }
