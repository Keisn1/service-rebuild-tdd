package mux_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Keisn1/note-taking-app/domain/web/auth"
	"github.com/Keisn1/note-taking-app/domain/web/mid"
	"github.com/Keisn1/note-taking-app/domain/web/mux"
	"github.com/Keisn1/note-taking-app/foundation/common"
	"github.com/Keisn1/note-taking-app/foundation/web"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Run("Single route example", func(t *testing.T) {
		cfg := mux.Config{}

		testRoutes := func(api *web.App, cfg mux.Config) {
			fetch := func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "Hello from fetch") }
			get := func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "Hello from get") }
			api.Handle("/fetch", http.HandlerFunc(fetch))
			api.Handle("/get", http.HandlerFunc(get))
		}

		api := mux.NewAPI(testRoutes, cfg)

		testCases := []struct {
			endpoint   string
			statusCode int
			want       string
		}{
			{endpoint: "/fetch", statusCode: http.StatusOK, want: "Hello from fetch"},
			{endpoint: "/get", statusCode: http.StatusOK, want: "Hello from get"},
		}
		for _, tc := range testCases {
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.endpoint, nil)
			api.ServeHTTP(resp, req)
			assert.Equal(t, tc.statusCode, resp.Code)
			assert.Equal(t, tc.want, resp.Body.String())
		}
	})

	t.Run("Example with authentication", func(t *testing.T) {
		key := common.MustGenerateRandomKey(32)
		cfg := mux.Config{Auth: auth.NewAuth(auth.MustNewJWTService(key))}

		testRoutes := func(api *web.App, cfg mux.Config) {
			authen := mid.Authenticate(cfg.Auth)
			fetch := func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "Hello from fetch") }
			api.Handle("/fetch", authen(http.HandlerFunc(fetch)))
		}

		api := mux.NewAPI(testRoutes, cfg)

		testCases := []struct {
			setupHeader func(*http.Request)
			statusCode  int
			want        string
		}{
			{setupHeader: func(r *http.Request) {}, statusCode: http.StatusForbidden, want: "failed authentication"},
			{setupHeader: func(r *http.Request) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
				tokenS, err := token.SignedString(key)
				assert.NoError(t, err)
				r.Header.Set("Authorization", "Bearer "+tokenS)
			},
				statusCode: http.StatusOK,
				want:       "Hello from fetch",
			},
		}

		for _, tc := range testCases {
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/fetch", nil)
			tc.setupHeader(req)
			api.ServeHTTP(resp, req)
			assert.Equal(t, tc.statusCode, resp.Code)
			assert.Contains(t, resp.Body.String(), tc.want)
		}
	})
}
