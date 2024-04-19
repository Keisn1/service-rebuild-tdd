package authMid

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Keisn1/note-taking-app/app/auth"
	"github.com/go-chi/chi/v5"
)

type MidHandler func(http.Handler) http.Handler

func Authenticate(a auth.AuthInterface) MidHandler {
	m := func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "userID")
			bearerToken := r.Header.Get("Authorization")

			_, err := a.Authenticate(userID, bearerToken)
			if err != nil {
				http.Error(w, "Failed Authentication", http.StatusForbidden)
				slog.Info(
					fmt.Sprintf(
						"Failed Authentication userID \"%v\" bearerToken \"%v\": %v",
						userID, bearerToken, err),
				)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	}
	return m
}
