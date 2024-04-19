package web

import "net/http"

type MidHandler func(http.Handler) http.Handler
