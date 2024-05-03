package note

import (
	"github.com/google/uuid"
)

type Note struct {
	NoteID  uuid.UUID
	Title   Title
	Content Content
	UserID  uuid.UUID
}

type UpdateNote struct {
	Title   Title
	Content Content
	UserID  uuid.UUID
}

type Content struct {
	content *string
}

type Title struct {
	title *string
}

func (n *Note) GetID() uuid.UUID { return n.NoteID }

func (n *Note) SetID(id uuid.UUID) { n.NoteID = id }

func (n *Note) GetTitle() Title       { return n.Title }
func (n *UpdateNote) GetTitle() Title { return n.Title }

func (n *Note) SetTitle(title string) { n.Title.Set(title) }

func (n *Note) GetContent() Content       { return n.Content }
func (n *UpdateNote) GetContent() Content { return n.Content }

func (n *Note) SetContent(content string) { n.Content.Set(content) }

func (n *Note) GetUserID() uuid.UUID { return n.UserID }

func NewTitle(title string) Title {
	return Title{title: &title}
}

func (tt Title) Set(title string) {
	*tt.title = title
}

func (tt Title) String() string {
	if tt.IsEmpty() {
		return ""
	}
	return *tt.title
}

func (tt Title) IsEmpty() bool { return tt.title == nil }

func NewContent(content string) Content {
	return Content{content: &content}
}

func (c Content) Set(content string) {
	*c.content = content
}

func (c Content) String() string {
	if c.IsEmpty() {
		return ""
	}
	return *c.content
}

func (c Content) IsEmpty() bool { return c.content == nil }
