package usernote

import (
	ents "github.com/Keisn1/note-taking-app/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetTitle(t *testing.T) {
	u := NewUserNote("title", "", uuid.New())
	u.SetTitle("newTitle")
	assert.Equal(t, ents.Title("newTitle"), u.GetTitle())
}
