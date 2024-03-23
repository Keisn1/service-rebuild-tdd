package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"strings"

	"github.com/golang-jwt/jwt"
)

type JWT string

type Auth struct{}

func (a *Auth) getTokenString(r *http.Request) (string, error) {
	parts := strings.Split(r.Header.Get("Authorization"), " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("expected authorization header format: Bearer <token>")
	}
	return parts[1], nil
}

func (a *Auth) parseTokenString(tokenS string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenS, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		secret := []byte(os.Getenv("JWT_SECRET_KEY"))
		return secret, nil
	})
	return token, err
}

func JWTAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := &Auth{}
		tokenString, err := a.getTokenString(r)
		if err != nil {
			http.Error(w, "Failed Authorization", http.StatusForbidden)
			slog.Info("Failed Authorization: ", err)
		}

		_, err = a.parseTokenString(tokenString)

		if err != nil {
			http.Error(w, "Failed Authorization", http.StatusForbidden)
			slog.Info("Failed Authorization: Token invalid", err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
}
