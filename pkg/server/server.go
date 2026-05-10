package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	shared_domain "github.com/hernanhrm/budget-forge/pkg/shared_domain"
)

type Config struct {
	Port               int
	ServiceName        string
	Debug              bool
	CORSAllowedOrigins []string
}

type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}

type Server struct {
	config   Config
	logger   shared_domain.Logger
	checkers map[string]HealthChecker

	Handler http.Handler
	srv     *http.Server
}

func NewServer(config Config, log shared_domain.Logger, routeSetupFunc func(*http.ServeMux)) *Server {
	s := &Server{
		config:   config,
		logger:   log.With("component", "server"),
		checkers: make(map[string]HealthChecker),
	}

	mux := http.NewServeMux()

	s.setupBaseRoutes(mux)

	if routeSetupFunc != nil {
		routeSetupFunc(mux)
	}

	s.Handler = withMiddleware(mux, s.config, s.logger)

	s.logger.Info("HTTP server initialized",
		"port", config.Port,
		"service", config.ServiceName,
	)

	return s
}

func (s Server) RegisterHealthChecker(name string, checker HealthChecker) {
	s.checkers[name] = checker
}

func (s Server) setupBaseRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /ping", s.handlePing)
}

func (s Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.config.Port)

	s.srv = &http.Server{
		Addr:         addr,
		Handler:      s.Handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.logger.WithContext(ctx).Info("starting HTTP server",
		"address", addr,
		"service", s.config.ServiceName,
	)

	errCh := make(chan error, 1)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("server: %w", err)
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		s.logger.WithContext(ctx).Info("shutting down HTTP server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return s.srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func (s Server) Shutdown(ctx context.Context) error {
	s.logger.WithContext(ctx).Info("shutting down HTTP server")
	return s.srv.Shutdown(ctx)
}
