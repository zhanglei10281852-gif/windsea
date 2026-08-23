package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zhanglei10281852-gif/windsea/internal/config"
	"github.com/zhanglei10281852-gif/windsea/internal/db"
	"github.com/zhanglei10281852-gif/windsea/internal/httpapi"
	"github.com/zhanglei10281852-gif/windsea/internal/observability"
	"github.com/zhanglei10281852-gif/windsea/internal/service"
	"github.com/zhanglei10281852-gif/windsea/internal/worker"
)

func main() {
	cfg := config.Load()
	logger := observability.NewLogger(os.Stdout, cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	database, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database open failed", "error", err)
		os.Exit(1)
	}
	defer database.Close()
	services := service.NewRegistry(database, logger, cfg.Clock)
	workers := worker.NewSupervisor(services, logger, cfg.WorkerInterval)
	workers.Start(ctx)
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: httpapi.NewRouter(services, logger), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		logger.Info("windsea listening", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", "error", err)
			stop()
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
	workers.Stop()
}
