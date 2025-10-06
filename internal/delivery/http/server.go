package http

import (
	"context"
	"fmt"
	"github.com/ttrtcixy/workout/internal/config"
	"github.com/ttrtcixy/workout/internal/delivery/http/handlers"
	"github.com/ttrtcixy/workout/internal/logger"
	"net/http"
)

// todo panic recovery

type Server struct {
	cfg *config.HttpServer
	log logger.Logger
	srv *http.Server
}

func New(cfg *config.HttpServer, log logger.Logger, handlers *handlers.Handlers) *Server {
	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: NewRouter(handlers).Handler(),
	}

	return &Server{
		cfg: cfg,
		log: log,
		srv: srv,
	}
}

// Start http server
func (s *Server) Start(ctx context.Context) error {
	s.log.Info("[+] starting http server on: %s", s.cfg.Addr())

	return s.srv.ListenAndServe()
}

// Close http server
func (s *Server) Close(ctx context.Context) error {
	const op = "http.server.close"
	// todo
	//ctx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout())
	//defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
