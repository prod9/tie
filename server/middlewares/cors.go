package middlewares

import (
	"github.com/rs/cors"
	"net/http"
	"tie.prodigy9.co/config"
)

func AddCorsAllowAll(handler http.Handler, cfg *config.Config) http.Handler {
	return cors.AllowAll().Handler(handler)
}
