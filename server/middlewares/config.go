package middlewares

import (
	"net/http"
	"tie.prodigy9.co/config"
)

func AddConfig(handler http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(resp, config.NewRequest(r, cfg))
	})
}
