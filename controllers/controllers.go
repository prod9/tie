package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	ErrNotFound     = &errImpl{"not_found", "not found"}
	ErrUnauthorized = &errImpl{"unauthorized", "unauthorized"}
	ErrInternal     = &errImpl{"internal", "internal server error"}
	ErrPermission   = &errImpl{"permission", "don't have permission to access"}
	ErrBadrequest   = &errImpl{"bad_request", "bad request"}

	allControllers = []Interface{
		HomeController{},
		CacheController{},
	}
)

type errImpl struct {
	code    string
	message string
}

func (i *errImpl) Code() string  { return i.code }
func (i *errImpl) Error() string { return i.message }

type Interface interface {
	Mount(router *httprouter.Router) error
}

func MountAll(router *httprouter.Router) error {
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
