package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWT interface {
	Verify(tokenS string) (Claims, error)
}

type jwtToken struct {
	key string
}

func NewJWT(key string) *jwtToken {
	return &jwtToken{
		key: key,
	}
}

type Claims jwt.MapClaims

type JWTPayload struct {
	jwt.RegisteredClaims
	Token string
}

func (j *jwtToken) CreateToken(userID uuid.UUID) *JWTPayload {
	payload := &JWTPayload{}
	payload.Subject = userID.String()

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	// tokenS, _ := token.SignedString([]byte(j.key))

	// payload.Token = tokenS
	return payload
}

func (j *jwtToken) Verify(tokenS string) (Claims, error) {
	claims, err := j.parseTokenString(tokenS)
	if err != nil {
		return nil, fmt.Errorf("verify: %w", err)
	}

	// if err := j.checkExpSet(claims); err != nil {
	// 	return nil, fmt.Errorf("verify: %w", err)
	// }

	// if err := j.checkIssuer(claims); err != nil {
	// 	return nil, fmt.Errorf("verify: %w", err)
	// }

	// if err := j.checkSubject(userID, claims); err != nil {
	// 	return nil, fmt.Errorf("verify: %w", err)
	// }
	return claims, err
}

func (j *jwtToken) parseTokenString(tokenS string) (Claims, error) {
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

func (j *jwtToken) keyFunc(token *jwt.Token) (interface{}, error) {
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
