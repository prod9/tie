package server

import (
	"net/http"
	"tie.prodigy9.co/config"
	"tie.prodigy9.co/server/middlewares"
)

type Middleware func(http.Handler, *config.Config) http.Handler

var stack = []Middleware{
	middlewares.AddConfig,
	middlewares.AddRequestLogging,
	middlewares.AddCorsAllowAll,
	middlewares.AddDataContext,
	middlewares.AddControllers,
}

type Server struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{cfg}
}

func (s *Server) Start() error {
	srv := http.NotFoundHandler()
	for i := len(stack) - 1; i >= 0; i-- {
		srv = stack[i](srv, s.cfg)
	}

	return http.ListenAndServe(s.cfg.ListenAddr(), srv)
}
