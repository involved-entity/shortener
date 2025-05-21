package internal

import (
	"shortener/internal/api"
	"shortener/internal/api/urls"
	"shortener/internal/api/users"
	"shortener/internal/database"
	"shortener/internal/machinery"
	"shortener/internal/redis"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(config *Config) {
	log := SetupLogger(config.Env)

	database.Init(config.DSN)
	machinery.Init(config.Mail.Email, config.Mail.Password)
	redis.Init()
	redisClient := redis.GetClient()

	defer redisClient.Close()

	log.Info("Starting shortener service", "env", config.Env)

	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/api/register", users.Register)
	e.POST("/api/login", users.Login(config.TTL, config.SECRET))

	e.GET("/_/:shortCode", urls.GetURL)

	authProtected := e.Group("")
	authProtected.Use(echojwt.WithConfig(echojwt.Config{SigningKey: []byte(config.SECRET), Skipper: api.OptionalJWT(
		[]string{"/api/urls"},
	)}))

	authProtected.POST("/api/urls", urls.SaveURL)
	authProtected.DELETE("/api/urls/:shortCode", urls.DeleteURL)

	e.Start(config.Address)
}
