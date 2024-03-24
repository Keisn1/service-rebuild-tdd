package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"strings"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
)

type key int

const userIDKey key = 1

type AuthInterface interface {
	Authenticate(userID, bearerToken string) (jwt.Claims, error)
}

type Auth struct{}

func (a *Auth) getTokenString(bearerToken string) (string, error) {
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("expected authorization header format: Bearer <token>")
	}
	return parts[1], nil
}

func (a *Auth) parseTokenString(tokenS string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenS, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		secret := []byte(os.Getenv("JWT_SECRET_KEY"))
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing tokenString: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("error extracting claims")
	}
}

func (a *Auth) isUserEnabled(userID string, claims jwt.MapClaims) error {
	if userID != claims["sub"] {
		return errors.New("user not enabled")
	}
	return nil
}

func (a *Auth) Authenticate(userID, bearerToken string) (jwt.Claims, error) {
	tokenS, err := a.getTokenString(bearerToken)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	claims, err := a.parseTokenString(tokenS)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	if err := a.isUserEnabled(userID, claims); err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	return claims, nil
}

type MidHandler func(http.Handler) http.Handler

func NewJwtMidHandler(a AuthInterface) MidHandler {
	m := func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "userID")
			bearerToken := r.Header.Get("Authorization")

			_, err := a.Authenticate(userID, bearerToken)
			if err != nil {
				http.Error(w, "Failed Authentication", http.StatusForbidden)
				slog.Info("Failed Authentication: ", err)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	}
	return m
}

func main() {
}
