package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"smartway-test/internal/config"
	"smartway-test/internal/storage/postgresql"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	//config
	cfg := config.New()
	//log
	log := setupLogger(cfg.Env)
	log.Info(
		"start api",
		slog.String("env", cfg.Env),
		slog.String("version", "1"),
	)
	log.Debug("debug messages are enabled")
	//storage
	storage, err := postgresql.New(cfg.ConnectionString)
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}
	_ = storage
	//router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api", func(r chi.Router) {
		//r.Get("/tickets")
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	//Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", err)
		}
	}()
	log.Info("server started")

	<-stop
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", err)

		os.Exit(1)
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
