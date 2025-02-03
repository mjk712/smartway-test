package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"smartway-test/internal/config"
	http_server "smartway-test/internal/http-server"
	"smartway-test/internal/service"
	"smartway-test/internal/storage"
	"smartway-test/internal/tools"
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
	//ctx
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//storage
	repo, err := storage.New(cfg.ConnectionString)
	if err != nil {
		log.Error("failed to init storage", tools.ErrAttr(err))
		os.Exit(1)
	}
	//services
	flightService := service.NewFlightService(repo)

	//server
	serv := http_server.NewServer(mainCtx, log, cfg, flightService)

	log.Info("starting server", slog.String("address", cfg.Address))

	//Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := serv.ListenAndServe(); err != nil {
			log.Error("failed to start server", tools.ErrAttr(err))
		}
	}()
	log.Info("server started")

	<-stop
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := serv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", tools.ErrAttr(err))

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
