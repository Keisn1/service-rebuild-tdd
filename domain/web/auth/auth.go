package auth

import (
	"errors"
	"fmt"
	"strings"
)

type AuthInterface interface {
	Authenticate(bearerToken string) (*Claims, error)
}

type Auth struct {
	jwtSvc JWTService
}

func NewAuth(jwtS JWTService) Auth {
	return Auth{jwtSvc: jwtS}
}

func (a Auth) Authenticate(bearerToken string) (*Claims, error) {
	tokenS, err := getBearerToken(bearerToken)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	claims, err := a.jwtSvc.Verify(tokenS)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	return claims, nil
}

func getBearerToken(bearerToken string) (string, error) {
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("expected authorization header format: Bearer <token>")
	}
	return parts[1], nil
}
