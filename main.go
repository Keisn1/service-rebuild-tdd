package main

import (
	"net/http"

	"github.com/Keisn1/note-taking-app/server"
	"log"
)

func main() {
	handler := http.HandlerFunc(server.NotesService)
	log.Fatal(http.ListenAndServe(":3000", handler))
}
