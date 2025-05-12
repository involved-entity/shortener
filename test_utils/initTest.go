package testutils

import (
	conf "shortener/internal/config"
	"shortener/internal/database"
	"shortener/internal/machinery"
	"shortener/internal/redis"

	r "github.com/go-redis/redis/v9"
)

func InitTest() *r.Client {
	conf.MustLoad()
	config := conf.GetConfig()

	database.Init(config.DSN)
	machinery.Init(config.Mail.Email, config.Mail.Password, config.Machinery.Broker, config.Machinery.ResultBackend)
	redis.Init(config.Redis.Address, config.Redis.Password, config.Redis.DB)
	return redis.GetClient()
}
