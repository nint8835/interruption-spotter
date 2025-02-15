package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/benbjohnson/hashfs"

	"github.com/nint8835/interruption-spotter/pkg/config"
	"github.com/nint8835/interruption-spotter/pkg/database"
	"github.com/nint8835/interruption-spotter/pkg/server/static"
)

type Server struct {
	logger      *slog.Logger
	cfg         *config.Config
	db          *database.Queries
	srv         *http.Server
	mux         *http.ServeMux
	stoppedChan chan struct{}
}

func (s *Server) run() {
	defer close(s.stoppedChan)

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("Failed to serve site", "err", err)
		return
	}
}

func (s *Server) Start() {
	s.logger.Debug("Starting server")
	go s.run()
}

func (s *Server) Stop(ctx context.Context) {
	s.logger.Debug("Stopping server")

	err := s.srv.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", "err", err)
	}
	<-s.stoppedChan

	s.logger.Debug("Server stopped")
}

func New(cfg *config.Config, db *database.Queries) *Server {
	mux := http.NewServeMux()

	instance := &Server{
		logger: slog.Default().With("component", "server"),
		cfg:    cfg,
		db:     db,
		mux:    mux,
		srv: &http.Server{
			Addr:    cfg.BindAddr,
			Handler: mux,
		},
		stoppedChan: make(chan struct{}),
	}

	mux.Handle("GET /static/", http.StripPrefix("/static/", hashfs.FileServer(static.HashFS)))

	mux.HandleFunc("GET /feed", instance.handleFeed)
	mux.HandleFunc("GET /{$}", instance.handleIndex)

	return instance
}
