package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ibrahimijai/go-http-server/config"
	"github.com/ibrahimijai/go-http-server/models"
	"github.com/ibrahimijai/go-http-server/routes"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

// run lets main stay tiny and return errors normally instead of calling
// os.Exit everywhere (os.Exit skips deferred calls).
func run() error {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	var st models.URLStore
	if cfg.DatabaseURL != "" {
		pg, err := models.NewPostgres(context.Background(), cfg.DatabaseURL)
		if err != nil {
			return fmt.Errorf("connecting to postgres: %w", err)
		}
		defer pg.Close()
		logger.Info("using postgres store")
		st = pg
	} else {
		logger.Info("DATABASE_URL not set, using in-memory store")
		st = models.NewMemory()
	}

	router := routes.New(logger, st)

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("server starting", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		logger.Info("shutdown signal received, draining connections...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	logger.Info("server stopped cleanly")
	return nil
}
