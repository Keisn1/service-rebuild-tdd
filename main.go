package main

import (
	"net/http"
)

func JWTAuthenticationMiddleware(next http.Handler) http.Handler {
	return next
}

func ValidateToken(t string) (claims map[string]string, err error) {
	claims = map[string]string{
		"iss": "note-taking-app",
	}
	return claims, nil
}

func main() {

}
