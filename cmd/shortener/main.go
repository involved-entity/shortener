package main

import (
	"shortener/internal"
	"shortener/internal/database"
	customMW "shortener/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config := internal.MustLoad()
	log := internal.SetupLogger(config.Env)

	log.Info("Starting shortener service", "env", config.Env)
	log.Info("HTTP server address", "address", config.HTTPServer.Address)

	database.Init(config.DSN)
	connection := database.GetDB()

	log.Info("Database connection", "connection", connection)

	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(customMW.NewLogMiddleware(log))
}
