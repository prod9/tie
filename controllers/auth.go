package controllers

import (
	"net/http"
	"strings"
	"tie.prodigy9.co/config"
	"tie.prodigy9.co/controllers/render"
)

func RequireAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			render.Error(resp, req, http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		cfg, token := config.FromRequest(req), auth[len("Bearer "):]
		if strings.ToLower(token) != strings.ToLower(cfg.AdminToken()) {
			render.Error(resp, req, http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		handler.ServeHTTP(resp, req)
	})
}
