package auth

import (
	"errors"
	"fmt"
	"os"

	"strings"

	"github.com/golang-jwt/jwt/v5"
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

func (a *Auth) checkIssuer(claims jwt.MapClaims) error {
	issuer := os.Getenv("JWT_NOTES_ISSUER")
	if issuer != claims["iss"] {
		return errors.New("incorrect Issuer")
	}
	return nil
}

func (a *Auth) checkExpSet(claims jwt.MapClaims) error {
	if _, ok := claims["exp"]; !ok {
		return fmt.Errorf("authenticate: no expiration date set")
	}
	return nil
}

func (a *Auth) Authenticate(userID, bearerToken string) error {
	tokenS, err := a.getTokenString(bearerToken)
	if err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}

	claims, err := a.parseTokenString(tokenS)
	if err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}

	if err := a.checkExpSet(claims); err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}

	if err := a.checkIssuer(claims); err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}

	if err := a.isUserEnabled(userID, claims); err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}
	return nil
}
