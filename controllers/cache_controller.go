package controllers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"tie.prodigy9.co/config"
	"tie.prodigy9.co/controllers/render"
	"tie.prodigy9.co/domain"
	"tie.prodigy9.co/internal/cache"
	"time"
)

type CacheController struct{}

var (
	_         Interface = CacheController{}
	tiesCache cache.Interface[[]*domain.Tie]
)

func init() {
	// TODO: Configure on mount?
	cfg := config.MustConfigure()
	if cfg.RedisURL() == "" {
		cfg.Println("using in-memory cache")
		tiesCache = cache.Basic[[]*domain.Tie]()
	} else {
		cfg.Println("using redis cache")
		tiesCache = cache.Redis[[]*domain.Tie](cfg, "tie-ties")
	}
}

func (c CacheController) Mount(router chi.Router) error {
	router.Route("/", func(r chi.Router) {
		r.Use(RequireAuth)
		r.Delete("/", c.Invalidate)
	})
	return nil
}

func (c CacheController) Invalidate(resp http.ResponseWriter, req *http.Request) {
	cfg, ctx := config.FromRequest(req), req.Context()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	go func() {
		defer cancel()
		if err := tiesCache.Invalidate(ctx); err != nil {
			cfg.Println("problem invalidating cache:", err)
		}
	}()

	render.JSON(resp, req, domain.CurrentStatus())
}
