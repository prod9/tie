package controllers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"tie.prodigy9.co/controllers/render"
	"tie.prodigy9.co/domain"
)

type HomeController struct{}

var _ Interface = HomeController{}

func (c HomeController) Mount(router *httprouter.Router) error {
	router.GET("/", c.Index)
	return nil
}

func (c HomeController) Index(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	render.JSON(resp, req, domain.CurrentStatus())
}
