package server

import (
	"encoding/json"
	"net/http"
)

func NotesService(w http.ResponseWriter, r *http.Request) {
	notes := []string{"Note number 1", "Note number 2"}
	json.NewEncoder(w).Encode(notes)
}
