package web

import "net/http"

type App struct {
	mux *http.ServeMux
}

func NewApp() *App {
	return &App{mux: http.NewServeMux()}
}

func (a *App) Handle(path string, handler http.Handler) {
	a.mux.Handle(path, handler)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
