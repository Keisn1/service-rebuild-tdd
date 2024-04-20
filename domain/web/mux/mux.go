package mux

import (
	"net/http"

	"github.com/Keisn1/note-taking-app/domain/web/auth"
	"github.com/Keisn1/note-taking-app/foundation/web"
)

type Config struct {
	Auth auth.Auth
}

type RouteAdder func(api *web.App, cfg Config)

func NewAPI(add RouteAdder, cfg Config) http.Handler {
	app := web.NewApp()
	add(app, cfg)
	return app
}
