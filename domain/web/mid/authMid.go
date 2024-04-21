package mid

import (
	"log/slog"
	"net/http"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/Keisn1/note-taking-app/domain/web/auth"
	"github.com/Keisn1/note-taking-app/foundation/web"
	"github.com/google/uuid"
)

type contextUserIDKey int

const UserIDKey contextUserIDKey = 1

func AuthorizeNote(ns note.NotesService) web.MidHandler {
	m := func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			noteID, _ := uuid.Parse(r.PathValue("note_id"))
			userID := r.Context().Value(UserIDKey).(uuid.UUID)
			n, _ := ns.GetNoteByID(noteID)
			if n.GetUserID() != userID {
				http.Error(w, "", http.StatusForbidden)
				return
			}
		}
		return http.HandlerFunc(h)
	}
	return m
}

func Authenticate(a auth.AuthInterface) web.MidHandler {
	m := func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			bearerToken := r.Header.Get("Authorization")
			_, err := a.Authenticate(bearerToken)
			if err != nil {
				http.Error(w, "failed authentication", http.StatusForbidden)
				slog.Info("failed authentication")
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	}
	return m
}
