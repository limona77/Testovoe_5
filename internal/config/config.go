package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		HTTP
		PG
		Env string `env:"ENV"`
	}
	HTTP struct {
		Port string `env:"HTTP_PORT"`
	}
	PG struct {
		URL string ` env:"PG_URL_LOCALHOST"`
	}
)

func NewConfig() *Config {
	cfg := &Config{}
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
