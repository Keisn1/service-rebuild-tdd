package mid

import (
	"log/slog"
	"net/http"

	"github.com/Keisn1/note-taking-app/domain/web/auth"
	"github.com/Keisn1/note-taking-app/foundation/web"
)

func Authenticate(a auth.AuthInterface) web.MidHandler {
	m := func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			bearerToken := r.Header.Get("Authorization")

			_, err := a.Authenticate(bearerToken)
			if err != nil {
				http.Error(w, "Failed Authentication", http.StatusForbidden)
				slog.Info("Failed Authentication")
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	}
	return m
}
