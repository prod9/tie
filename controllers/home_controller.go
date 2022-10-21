package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"tie.prodigy9.co/config"
	"tie.prodigy9.co/controllers/render"
	"tie.prodigy9.co/domain"
)

type HomeController struct{}

var _ Interface = HomeController{}

func (c HomeController) Mount(cfg *config.Config, router chi.Router) error {
	router.Get("/__health", c.Index)
	router.Get("/__smoke", c.Smoke)
	return nil
}

func (c HomeController) Index(resp http.ResponseWriter, req *http.Request) {
	render.JSON(resp, req, domain.CurrentStatus())
}
func (c HomeController) Smoke(resp http.ResponseWriter, req *http.Request) {
	// for smoke testing
	m := map[string]any{"message": "Hello Smoke Tests"}
	for key, values := range req.URL.Query() {
		m[key] = fmt.Sprint(values[0])
	}

	render.JSON(resp, req, m)
}
