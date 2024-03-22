package main

import (
	"net/http"
)

type JWTToken string

func JWTAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val, _ := r.Context().Value(JWTToken("token")).(bool)
		if !val {
			w.WriteHeader(http.StatusForbidden)
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
