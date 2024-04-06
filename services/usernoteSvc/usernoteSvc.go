// Package usernoteSvc wraps the data/store layer
// handles Crud operations on aggregate usernote
// make changes persistent by calling data/store layer
package usernoteSvc

type Note struct {
	Owner    string
	NoteName string
	NoteText string
}

func GetNoteByName(name string) Note {
	return Note{
		Owner:    "rob",
		NoteName: "robs note",
		NoteText: "robs note text",
	}
}
