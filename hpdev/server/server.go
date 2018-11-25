package server

import (
	"net/http"
)

// New -
func New(cfg Config) (*Server, error) {
	s := &http.Server{
		Handler: cfg.Handler,
		Addr:    cfg.Addr,
	}

	return &Server{
		Config: cfg,
		s:      s,
	}, nil
}

// Must -
func Must(cfg Config) *Server {
	s, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return s
}

// Config -
type Config struct {
	Addr    string
	Handler http.Handler
}

// Server -
type Server struct {
	Config
	s *http.Server
}

// Run -
func (s *Server) Run() error {
	return s.s.ListenAndServe()
}
