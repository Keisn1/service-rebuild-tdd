package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetNotes(t *testing.T) {
	t.Run("Server returns array of notes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/notes", nil)
		response := httptest.NewRecorder()
		NotesService(response, request)

		var got []string
		json.NewDecoder(response.Body).Decode(&got)
		want := []string{"Note number 1", "Note number 2"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf(`got = %v; want %v`, got, want)
		}
	})
}
