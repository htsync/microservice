package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func LoggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logger.Info("Request processed",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Duration("duration", duration))
		})
	}
}
