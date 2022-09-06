package controllers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"tie.prodigy9.co/controllers/render"
	"tie.prodigy9.co/domain"
)

type HomeController struct{}

var _ Interface = HomeController{}

func (c HomeController) Mount(router chi.Router) error {
	router.Get("/__health", c.Index)
	router.Get("/__smoke", c.Smoke)
	return nil
}

func (c HomeController) Index(resp http.ResponseWriter, req *http.Request) {
	render.JSON(resp, req, domain.CurrentStatus())
}
func (c HomeController) Smoke(resp http.ResponseWriter, req *http.Request) {
	// for smoke testing
	render.JSON(resp, req, "Hello Smoke Tests")
}
