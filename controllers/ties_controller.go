package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"tie.prodigy9.co/controllers/render"
	"tie.prodigy9.co/domain"
)

type TiesController struct{}

var _ Interface = TiesController{}

func (c TiesController) Mount(router chi.Router) error {
	router.Route("/ties", func(r chi.Router) {
		r.Use(RequireAuth)
		r.Get("/", c.Index)
		r.Post("/", c.Create)
		r.Delete("/{slug:[a-zA-Z0-9]+}", c.Delete)
	})

	router.Get("/{slug:[a-zA-Z0-9]+}", c.Redirect)
	return nil
}

func (c TiesController) Redirect(resp http.ResponseWriter, req *http.Request) {
	slug, tie := chi.URLParam(req, "slug"), &domain.Tie{}
	if err := domain.GetTieBySlug(req.Context(), tie, slug); err != nil {
		render.Error(resp, req, http.StatusBadRequest, err)
	} else {
		http.Redirect(resp, req, tie.TargetURL, http.StatusTemporaryRedirect)
	}
}

func (c TiesController) Index(resp http.ResponseWriter, req *http.Request) {
	ties := domain.NewList[*domain.Tie](nil)
	if err := domain.ListAllTies(req.Context(), ties); err != nil {
		render.Error(resp, req, http.StatusInternalServerError, err)
	} else {
		render.JSON(resp, req, ties)
	}
}

func (c TiesController) Create(resp http.ResponseWriter, req *http.Request) {
	action, tie := &domain.CreateTie{}, &domain.Tie{}
	if err := json.NewDecoder(req.Body).Decode(action); err != nil {
		render.Error(resp, req, http.StatusBadRequest, ErrBadRequest)
	} else if err = action.Validate(); err != nil {
		render.Error(resp, req, http.StatusBadRequest, err)
	} else if err = action.Execute(req.Context(), tie); err != nil {
		render.Error(resp, req, http.StatusInternalServerError, err)
	} else {
		render.JSON(resp, req, tie)
	}
}

func (c TiesController) Delete(resp http.ResponseWriter, req *http.Request) {
	slug := chi.URLParam(req, "slug")

	action, tie := &domain.DeleteTie{Slug: slug}, &domain.Tie{}
	if err := action.Validate(); err != nil {
		render.Error(resp, req, http.StatusBadRequest, err)
	} else if err := action.Execute(req.Context(), tie); err != nil {
		render.Error(resp, req, http.StatusInternalServerError, err)
	} else {
		render.JSON(resp, req, tie)
	}
}
