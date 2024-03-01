package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubNotesStore struct {
	notes []string
}

func (sns *StubNotesStore) GetAllNotes() []string {
	return sns.notes
}

func TestGetNotes(t *testing.T) {
	notesStore := StubNotesStore{
		notes: []string{"Note number 1", "Note number 2"},
	}
	notesServer := &NotesServer{NotesStore: &notesStore}

	t.Run("Server returns all Notes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/notes", nil)
		response := httptest.NewRecorder()
		notesServer.ServeHTTP(response, request)

		var got []string
		json.NewDecoder(response.Body).Decode(&got)
		want := []string{"Note number 1", "Note number 2"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf(`got = %v; want %v`, got, want)
		}
	})
}
