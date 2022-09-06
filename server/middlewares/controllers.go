package middlewares

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
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
			router, err = buildRouter()
		})

		if err != nil {
			render.Error(resp, req, 500, fmt.Errorf("router problems: %w", err))
		} else {
			router.ServeHTTP(resp, req)
		}
	})
}

func buildRouter() (http.Handler, error) {
	router := httprouter.New()
	if err := controllers.MountAll(router); err != nil {
		return nil, err
	}

	router.NotFound = http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		render.Error(resp, r, 404, controllers.ErrNotFound)
	})
	return router, nil
}
