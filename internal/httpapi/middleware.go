package httpapi

import (
	"github.com/zhanglei10281852-gif/windsea/internal/observability"
	"log/slog"
	"net/http"
)

func withRequestID(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = observability.NewRequestID()
		}
		w.Header().Set("X-Request-ID", id)
		ctx := observability.WithRequestID(r.Context(), id)
		logger.Debug("request", "method", r.Method, "path", r.URL.Path, "request_id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func withRecovery(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if value := recover(); value != nil {
				logger.Error("panic recovered", "panic", value)
				writeJSON(w, http.StatusInternalServerError, map[string]any{"code": "internal_error", "message": "internal server error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
