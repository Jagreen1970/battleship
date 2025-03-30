package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Jagreen1970/battleship/internal/app"
)

type Server struct {
	httpServer *http.Server
	cfg        app.ServerConfig
}

func New(cfg app.ServerConfig) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      http.DefaultServeMux,
		ReadTimeout:  s.cfg.Timeout,
		WriteTimeout: s.cfg.Timeout,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Timeout)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
