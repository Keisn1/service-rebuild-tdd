package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"strings"

	"github.com/golang-jwt/jwt"
)

type key int

const userIDKey key = 1

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

func (a *Auth) isUserEnabled(ctx context.Context, claims jwt.MapClaims) error {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return errors.New("error parsing userIDKey to key")
	}
	if userID != claims["sub"] {
		return errors.New("user not enabled")
	}
	return nil
}

func JWTAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := &Auth{}
		tokenString, err := a.getTokenString(r.Header.Get("Authorization"))
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
