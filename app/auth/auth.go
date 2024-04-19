package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Keisn1/note-taking-app/foundation/jwtSvc"
)

type key int

const userIDKey key = 1

type AuthInterface interface {
	Authenticate(userID, bearerToken string) (*jwtSvc.Claims, error)
}

type Auth struct {
	jwtS jwtSvc.JWTService
}

func NewAuth(jwtS jwtSvc.JWTService) Auth {
	return Auth{jwtS: jwtS}
}

func (a Auth) Authenticate(userID, bearerToken string) (*jwtSvc.Claims, error) {
	tokenS, err := getJWTTokenString(bearerToken)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	claims, err := a.jwtS.Verify(tokenS)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	return claims, nil
}

func getJWTTokenString(bearerToken string) (string, error) {
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("expected authorization header format: Bearer <token>")
	}
	return parts[1], nil
}
