package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/hernanhrm/budget-forge/pkg/httpresponse"
	shared_domain "github.com/hernanhrm/budget-forge/pkg/shared_domain"
)

func withMiddleware(mux *http.ServeMux, config Config, log shared_domain.Logger) http.Handler {
	var handler http.Handler = mux

	handler = recoveryMiddleware(handler)
	handler = loggingMiddleware(handler, log)
	handler = corsMiddleware(handler, config)

	return handler
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "error", rec)
				httpresponse.WriteInternalServerError(w, r, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler, log shared_domain.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)

		ctx := context.Background()
		log.WithContext(ctx).Info("request completed",
			"method", r.Method,
			"uri", r.RequestURI,
			"status", sw.status,
			"latency", time.Since(start),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(status int) {
	sw.status = status
	sw.ResponseWriter.WriteHeader(status)
}

func corsMiddleware(next http.Handler, config Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(config.CORSAllowedOrigins) > 0 {
			origin := r.Header.Get("Origin")
			for _, allowed := range config.CORSAllowedOrigins {
				if allowed == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					break
				}
			}
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, PUT, PATCH, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Datastar-Request")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
