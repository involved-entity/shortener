package main

import (
	"shortener/internal"
	"shortener/internal/database"
)

func main() {
	config := internal.MustLoad()
	log := internal.SetupLogger(config.Env)

	log.Info("Starting shortener service", "env", config.Env)
	log.Info("HTTP server address", "address", config.HTTPServer.Address)

	database.Init(config.DSN)
	connection := database.GetDB()

	log.Info("Database connection", "connection", connection)
}
