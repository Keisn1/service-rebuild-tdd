package note

import (
	"github.com/google/uuid"
)

type Note struct {
	noteID  uuid.UUID
	title   Title
	content Content
	userID  uuid.UUID
}

type NewNote struct {
	title   Title
	content Content
	userID  uuid.UUID
}

func MakeNote(noteID uuid.UUID, title Title, content Content, userID uuid.UUID) Note {
	return Note{
		noteID:  noteID,
		title:   title,
		content: content,
		userID:  userID,
	}
}

func MakeNoteFromNewNote(nN NewNote) Note {
	return Note{
		noteID:  uuid.New(),
		title:   nN.title,
		content: nN.content,
		userID:  nN.userID,
	}
}

func MakeNewNote(title Title, content Content, userID uuid.UUID) NewNote {
	return NewNote{
		title:   title,
		content: content,
		userID:  userID,
	}
}

func (n *Note) GetID() uuid.UUID { return n.noteID }

func (n *Note) SetID(id uuid.UUID) { n.noteID = id }

func (n *Note) GetTitle() Title { return n.title }

func (n *Note) SetTitle(title string) { n.title.Set(title) }

func (n *Note) GetContent() Content { return n.content }

func (n *Note) SetContent(content string) { n.content.Set(content) }

func (n *Note) GetUserID() uuid.UUID { return n.userID }

type Title struct {
	title *string
}

func NewTitle(title string) Title {
	return Title{title: &title}
}

func (tt Title) Set(title string) {
	*tt.title = title
}

func (tt Title) String() string {
	return *tt.title
}

func (tt Title) IsEmpty() bool { return tt.title == nil }

type Content struct {
	content *string
}

func NewContent(content string) Content {
	return Content{content: &content}
}

func (c Content) Set(content string) {
	*c.content = content
}

func (c Content) String() string {
	return *c.content
}

func (c Content) IsEmpty() bool { return c.content == nil }
