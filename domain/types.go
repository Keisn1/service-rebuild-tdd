package domain

type Note struct {
	NoteID int
	UserID int
	Note   string
}

type Notes []Note

type NotePost struct {
	Note string `json:"note"`
}
