package server

import (
	"encoding/json"
	"io"
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

func (sns *StubNotesStore) GetNotesById(id int) []string {
	if id == 1 {
		return []string{"Note 1 user 1", "Note 2 user 1"}
	} else {
		return []string{"Note 1 user 2", "Note 2 user 2"}
	}
}

func TestNotes(t *testing.T) {
	notesStore := StubNotesStore{
		notes: []string{"Note number 1", "Note number 2"},
	}
	notesServer := &NotesServer{NotesStore: &notesStore}

	t.Run("Server returns all Notes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/notes", nil)
		response := httptest.NewRecorder()
		notesServer.ServeHTTP(response, request)

		got := decodeJsonBody(response.Body)
		want := []string{"Note number 1", "Note number 2"}
		assertStringSliceEqual(t, got, want)
	})

	t.Run("Server returns all Notes for user 1", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/notes/1", nil)
		response := httptest.NewRecorder()
		notesServer.ServeHTTP(response, request)

		got := decodeJsonBody(response.Body)
		want := []string{"Note 1 user 1", "Note 2 user 1"}

		assertStringSliceEqual(t, got, want)
	})
	t.Run("Server returns all Notes for user 2", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/notes/2", nil)
		response := httptest.NewRecorder()
		notesServer.ServeHTTP(response, request)

		got := decodeJsonBody(response.Body)
		want := []string{"Note 1 user 2", "Note 2 user 2"}

		assertStringSliceEqual(t, got, want)
	})
}

func assertStringSliceEqual(t testing.TB, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func decodeJsonBody(r io.Reader) []string {
	var res []string
	json.NewDecoder(r).Decode(&res)
	return res
}
