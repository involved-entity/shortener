package main

import (
	"log/slog"
	"os"
	"shortener/internal"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	config := internal.MustLoad()
	log := setupLogger(config.Env)
	log.Info("Starting shortener service", "env", config.Env)
	log.Info("HTTP server address", "address", config.HTTPServer.Address)
	connection, error := internal.NewDatabase(config.DSN)
	log.Info("Database connection", "connection", connection, "error", error)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
