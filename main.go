package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

type JWTToken string

func JWTAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, _ := r.Context().Value(JWTToken("token")).(string)

		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			secret := []byte(os.Getenv("JWT_SECRET_KEY"))
			return secret, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("No valid JWTToken"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ValidateToken(t string) (claims map[string]string, err error) {
	claims = map[string]string{
		"iss": "note-taking-app",
	}
	return claims, nil
}

func main() {
}
