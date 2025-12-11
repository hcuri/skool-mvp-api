package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hcuri/skool-mvp-app/internal/config"
	"github.com/hcuri/skool-mvp-app/internal/db"
	apihttp "github.com/hcuri/skool-mvp-app/internal/http"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := config.Load()

	logger, err := initLogger(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	store := db.NewInMemoryStore()

	router := apihttp.NewRouter(store, logger)
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("starting server", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	}
	logger.Info("server stopped")
}

func initLogger(level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	if level != "" {
		var l zapcore.Level
		if err := l.Set(level); err == nil {
			cfg.Level = zap.NewAtomicLevelAt(l)
		}
	}

	return cfg.Build()
}
