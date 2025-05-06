package internal

import (
	"os"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"`
	DSN        string `yaml:"dsn" env:"DSN"`
	HTTPServer `yaml:"http_server"`
	JWT        `yaml:"jwt"`
}

type HTTPServer struct {
	Address      string `yaml:"address" env:"ADDRESS" env-default:"localhost:8080"`
	IddleTimeout int    `yaml:"iddle_timeout" env:"IDDLE_TIMEOUT" env-default:"60"`
	ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"60"`
	Timeout      int    `yaml:"timeout" env:"TIMEOUT" env-default:"60"`
}

type JWT struct {
	TTL    int    `yaml:"ttl" env:"TTL" env-default:"60"`
	SECRET string `yaml:"secret" env:"SECRET"`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "config/local.yml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not found")
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read config file: " + err.Error())
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config, func(c *mapstructure.DecoderConfig) {
		c.TagName = "yaml"
	}); err != nil {
		panic("failed to unmarshal config: " + err.Error())
	}

	return config
}
