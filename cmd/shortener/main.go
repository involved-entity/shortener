package main

import (
	"shortener/internal"
	"shortener/internal/api/urls"
	"shortener/internal/api/users"
	"shortener/internal/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := internal.MustLoad()
	log := internal.SetupLogger(config.Env)

	log.Info("Starting shortener service", "env", config.Env)
	log.Info("HTTP server address", "address", config.HTTPServer.Address)

	database.Init(config.DSN)
	connection := database.GetDB()

	log.Info("Database connection", "connection", connection)

	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/api/register", users.Register)
	e.POST("/api/login", users.Login(config.TTL, config.SECRET))

	e.GET("/_/:shortCode", urls.GetURL)
	e.POST("/api/urls", urls.SaveURL)
	e.DELETE("/api/urls/:shortCode", urls.DeleteURL)

	e.Start(config.Address)
}
