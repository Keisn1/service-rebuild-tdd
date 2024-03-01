package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

type StubNotesStore struct {
	notes        map[int][]string
	addNoteCalls []addNoteCall
}

type addNoteCall struct {
	userID int
	note   string
}

func (sns *StubNotesStore) AddNote(userID int, note string) {
	sns.addNoteCalls = append(sns.addNoteCalls, addNoteCall{userID, note})
}

func (sns *StubNotesStore) GetAllNotes() []string {
	var allNotes []string
	for _, notes := range sns.notes {
		allNotes = append(allNotes, notes...)
	}
	return allNotes
}

func (sns *StubNotesStore) GetNotesByID(id int) []string {
	return sns.notes[id]
}

func TestNotes(t *testing.T) {
	notesStore := StubNotesStore{
		notes: map[int][]string{
			1: {"Note 1 user 1", "Note 2 user 1"},
			2: {"Note 1 user 2", "Note 2 user 2"},
		},
	}
	notesServer := &NotesServer{NotesStore: &notesStore}

	t.Run("Server returns all Notes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/notes", nil)
		response := httptest.NewRecorder()
		notesServer.ServeHTTP(response, request)

		assertStatusCode(t, response.Result().StatusCode, http.StatusOK)

		got := decodeJsonBody(response.Body)
		want := []string{"Note 1 user 1", "Note 2 user 1", "Note 1 user 2", "Note 2 user 2"}
		assertStringSlicesAreEqual(t, got, want)
	})

	t.Run("Server returns all Notes for user 1", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/notes/1", nil)
		response := httptest.NewRecorder()
		notesServer.ServeHTTP(response, request)

		assertStatusCode(t, response.Result().StatusCode, http.StatusOK)

		got := decodeJsonBody(response.Body)
		want := []string{"Note 1 user 1", "Note 2 user 1"}
		assertSlicesHaveSameLength(t, got, want)
		assertStringSlicesAreEqual(t, got, want)
	})
	t.Run("Server returns all Notes for user 2", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/notes/2", nil)
		response := httptest.NewRecorder()
		notesServer.ServeHTTP(response, request)

		assertStatusCode(t, response.Result().StatusCode, http.StatusOK)

		got := decodeJsonBody(response.Body)
		want := []string{"Note 1 user 2", "Note 2 user 2"}
		assertSlicesHaveSameLength(t, got, want)
		assertStringSlicesAreEqual(t, got, want)
	})
	t.Run("Server returns zero Notes for user 100", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/notes/100", nil)
		response := httptest.NewRecorder()
		notesServer.ServeHTTP(response, request)

		assertStatusCode(t, response.Result().StatusCode, http.StatusNotFound)
	})
	t.Run("adds a note with POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/notes/1", nil)
		response := httptest.NewRecorder()
		notesServer.ServeHTTP(response, request)

		assertStatusCode(t, response.Result().StatusCode, http.StatusAccepted)
		assertLengthSlice(t, notesStore.addNoteCalls, 1)
	})
}

func assertLengthSlice[T any](t testing.TB, elements []T, want int) {
	t.Helper()
	if len(elements) != want {
		t.Errorf(`got = %v; want %v`, len(elements), want)
	}

}

func assertStringSlicesAreEqual(t testing.TB, got, want []string) {
	t.Helper()
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	sort.Slice(want, func(i, j int) bool { return want[i] < want[j] })
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func decodeJsonBody(r io.Reader) []string {
	var res []string
	json.NewDecoder(r).Decode(&res)
	return res
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf(`got = %v; want %v`, got, want)
	}
}

func assertSlicesHaveSameLength[T any](t testing.TB, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf(`len(got) = %v; len(want) %v`, len(got), len(want))
	}
}
