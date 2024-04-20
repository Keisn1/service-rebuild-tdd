package mux_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Keisn1/note-taking-app/domain/web/auth"
	"github.com/Keisn1/note-taking-app/domain/web/mid"
	"github.com/Keisn1/note-taking-app/domain/web/mux"
	"github.com/Keisn1/note-taking-app/foundation/common"
	"github.com/Keisn1/note-taking-app/foundation/web"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Run("Single route example", func(t *testing.T) {
		cfg := mux.Config{}
		dataFetch := "Hello from fetch"
		testRoutes := func(api *web.App, cfg mux.Config) {
			fetch := func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, dataFetch) }
			api.Handle("/fetch", http.HandlerFunc(fetch))
		}

		api := mux.NewAPI(testRoutes, cfg)

		ts := httptest.NewServer(api)
		req := httptest.NewRequest(http.MethodGet, ts.URL+"/fetch", nil)
		req.RequestURI = ""

		res, err := ts.Client().Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		respBodyBytes, _ := io.ReadAll(res.Body)
		if string(respBodyBytes) != dataFetch {
			t.Errorf(`got "%s", want "%s"`, respBodyBytes, dataFetch)
		}
	})

	t.Run("Example with authentication", func(t *testing.T) {
		key := common.MustGenerateRandomKey(32)
		cfg := mux.Config{Auth: auth.NewAuth(auth.MustNewJWTService(key))}

		testRoutes := func(api *web.App, cfg mux.Config) {
			authen := mid.Authenticate(cfg.Auth)
			fetch := func(w http.ResponseWriter, r *http.Request) {}
			api.Handle("/fetch", authen(http.HandlerFunc(fetch)))
		}

		api := mux.NewAPI(testRoutes, cfg)

		ts := httptest.NewServer(api)
		req := httptest.NewRequest(http.MethodGet, ts.URL+"/fetch", nil)
		req.RequestURI = ""

		res, err := ts.Client().Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, http.StatusForbidden, res.StatusCode)
	})
}
