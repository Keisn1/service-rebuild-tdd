package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Keisn1/note-taking-app/foundation/jwt"
)

type key int

const userIDKey key = 1

type UserStore interface {
	FindUserByID(userID string) error
}

type Auth struct {
	jwt jwt.JWT
}

func (a *Auth) Authenticate(userID, bearerToken string) (jwt.Claims, error) {
	tokenS, err := a.getJWTTokenString(bearerToken)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	claims, err := a.jwt.Verify(tokenS)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	return claims, nil
}

func (a *Auth) getJWTTokenString(bearerToken string) (string, error) {
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("expected authorization header format: Bearer <token>")
	}
	return parts[1], nil
}
