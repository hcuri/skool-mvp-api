package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hcuri/skool-mvp-app/internal/config"
	"github.com/hcuri/skool-mvp-app/internal/db"
	apihttp "github.com/hcuri/skool-mvp-app/internal/http"
)

func main() {
	cfg := config.Load()

	logger := log.New(os.Stdout, "api ", log.LstdFlags)
	store := db.NewInMemoryStore()

	router := apihttp.NewRouter(store, logger)
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Printf("starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("graceful shutdown failed: %v", err)
	}
}
