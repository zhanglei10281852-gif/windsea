package observability

import (
	"io"
	"log/slog"
)

func NewLogger(output io.Writer, level string) *slog.Logger {
	var configured slog.Level
	if level == "DEBUG" {
		configured = slog.LevelDebug
	} else if level == "WARN" {
		configured = slog.LevelWarn
	} else {
		configured = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: configured}))
}
