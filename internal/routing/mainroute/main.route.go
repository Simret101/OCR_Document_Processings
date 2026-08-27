package mainroute

import (
	"aidoc/platform/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []func(http.Handler) http.Handler
}

func RegisterRoute(grg chi.Router, routes []Route, log logger.Logger) {
	for _, route := range routes {
		var handler http.Handler = route.Handler
		for i := len(route.Middleware) - 1; i >= 0; i-- {
			handler = route.Middleware[i](handler)
		}
		grg.Method(route.Method, route.Path, handler)

	}
}
