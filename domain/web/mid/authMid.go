package mid

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Keisn1/note-taking-app/domain/core/note"
	"github.com/Keisn1/note-taking-app/domain/web/auth"
	"github.com/Keisn1/note-taking-app/foundation"
	"github.com/Keisn1/note-taking-app/foundation/web"
	"github.com/google/uuid"
)

func AuthorizeNote(ns note.Service) web.MidHandler {
	m := func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			noteID, err := uuid.Parse(r.PathValue("note_id"))
			if err != nil {
				http.Error(w, "", http.StatusForbidden)
				return
			}

			userID := r.Context().Value(foundation.UserIDKey).(uuid.UUID)
			n, err := ns.QueryByID(r.Context(), noteID)
			if err != nil {
				http.Error(w, "", http.StatusForbidden)
				return
			}

			if n.UserID != userID {
				http.Error(w, "", http.StatusForbidden)
				return
			}

			ctx := setNote(r.Context(), n)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	}
	return m
}

func Authenticate(a auth.AuthInterface) web.MidHandler {
	m := func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			bearerToken := r.Header.Get("Authorization")
			claims, err := a.Authenticate(bearerToken)
			if err != nil {
				http.Error(w, "failed authentication", http.StatusForbidden)
				slog.Info("failed authentication")
				return
			}

			userID, _ := uuid.Parse(claims.Subject)
			ctx := setUserID(r.Context(), userID)
			ctx = setClaims(ctx, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	}
	return m
}

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, foundation.UserIDKey, userID)
}

func GetUserID(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(foundation.UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}
	}
	return userID
}

func setNote(ctx context.Context, n note.Note) context.Context {
	return context.WithValue(ctx, foundation.NoteKey, n)
}

func GetNote(ctx context.Context) note.Note {
	n, ok := ctx.Value(foundation.NoteKey).(note.Note)
	if !ok {
		fmt.Println("Not ok")
		return note.Note{}
	}
	return n
}

func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, foundation.ClaimsKey, claims)
}

func GetClaims(ctx context.Context) auth.Claims {
	claims, ok := ctx.Value(foundation.ClaimsKey).(auth.Claims)
	if !ok {
		return auth.Claims{}
	}
	return claims
}
