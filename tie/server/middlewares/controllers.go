package middlewares

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sync"
	"tie.prodigy9.co/config"
	"tie.prodigy9.co/controllers"
	"tie.prodigy9.co/controllers/render"
)

func AddControllers(handler http.Handler, cfg *config.Config) http.Handler {
	var (
		once   sync.Once
		router http.Handler
		err    error
	)

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		once.Do(func() {
			router, err = buildRouter(cfg)
		})

		if err != nil {
			render.Error(resp, req, 500, fmt.Errorf("router problems: %w", err))
		} else {
			router.ServeHTTP(resp, req)
		}
	})
}

func buildRouter(cfg *config.Config) (http.Handler, error) {
	router := chi.NewRouter()
	if err := controllers.MountAll(cfg, router); err != nil {
		return nil, err
	}

	router.NotFound(func(resp http.ResponseWriter, r *http.Request) {
		render.Error(resp, r, 404, controllers.ErrNotFound)
	})
	return router, nil
}
