package api

type UrlParams struct {
	NoteID int
}

type NotePost struct {
	Note string `json:"note"`
}
