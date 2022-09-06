package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var (
	ErrNotFound     = &errImpl{"not_found", "not found"}
	ErrUnauthorized = &errImpl{"unauthorized", "unauthorized"}
	ErrInternal     = &errImpl{"internal", "internal server error"}
	ErrPermission   = &errImpl{"permission", "don't have permission to access"}
	ErrBadRequest   = &errImpl{"bad_request", "bad request"}

	allControllers = []Interface{
		HomeController{},
		CacheController{},
		TiesController{},
	}
)

type errImpl struct {
	code    string
	message string
}

func (i *errImpl) Code() string  { return i.code }
func (i *errImpl) Error() string { return i.message }

type Interface interface {
	Mount(router chi.Router) error
}

func MountAll(router chi.Router) error {
	for _, controller := range allControllers {
		if err := controller.Mount(router); err != nil {
			return err
		}
	}
	return nil
}

func ReadJSON(r *http.Request, obj interface{}) error {
	err := json.NewDecoder(r.Body).Decode(obj)
	if err != nil {
		return ErrInternal
	}
	return nil
}
