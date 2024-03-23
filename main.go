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

func getTokenString(r *http.Request) (string, error) {
	parts := strings.Split(r.Header.Get("Authorization"), " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("expected authorization header format: Bearer <token>")
	}
	return parts[1], nil
}

func parseTokenString(tokenS string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenS, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Info("unexpected signing method: %v", token.Header["alg"])
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
		tokenString, err := getTokenString(r)
		if err != nil {
			http.Error(w, "Failed Authorization", http.StatusForbidden)
			slog.Info("Failed Authorization: ", err)
		}

		_, err = parseTokenString(tokenString)

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
