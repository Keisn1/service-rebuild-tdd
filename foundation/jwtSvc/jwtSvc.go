package jwtSvc

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWT interface {
	Verify(tokenS string) (Claims, error)
}

type jwtSvc struct {
	key []byte
}

func NewJWT(key []byte) (*jwtSvc, error) {
	if len(key) < 32 {
		return nil, errors.New("key minLength 32")
	}
	return &jwtSvc{key: key}, nil
}

type Claims jwt.MapClaims

type JWTPayload struct {
	jwt.RegisteredClaims
	Token string
}

func (j *jwtSvc) CreateToken(userID uuid.UUID, d time.Duration) (*JWTPayload, error) {
	payload := &JWTPayload{}
	payload.Subject = userID.String()
	payload.ExpiresAt = jwt.NewNumericDate(time.Now().Add(d))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenS, err := token.SignedString(j.key)
	if err != nil {
		return nil, err
	}

	payload.Token = tokenS
	return payload, nil
}

func (j *jwtSvc) Verify(tokenS string) (*JWTPayload, error) {
	token, err := jwt.ParseWithClaims(tokenS, &JWTPayload{}, j.keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := token.Claims.(*JWTPayload)
	if !ok {
		return nil, errors.New("invalid token")
	}

	payload.Token = tokenS
	return payload, nil
}

func (j *jwtSvc) parseTokenString(tokenS string) (Claims, error) {
	token, err := jwt.Parse(tokenS, j.keyFunc)
	if err != nil {
		return nil, fmt.Errorf("error parsing tokenString: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return Claims(claims), nil
	} else {
		return nil, errors.New("error extracting claims")
	}
}

func (j *jwtSvc) keyFunc(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(j.key), nil
}

// func (j *JWT) checkSubject(userID string, claims Claims) error {
// 	if userID != claims["sub"] {
// 		return errors.New("invalid subject")
// 	}
// 	return nil
// }

// func (j *JWT) checkIssuer(claims Claims) error {
// 	issuer := os.Getenv("JWT_NOTES_ISSUER")
// 	if issuer != claims["iss"] {
// 		return errors.New("incorrect issuer")
// 	}
// 	return nil
// }

// func (j *JWT) checkExpSet(claims Claims) error {
// 	if _, ok := claims["exp"]; !ok {
// 		return fmt.Errorf("no expiration date set")
// 	}
// 	return nil
// }
