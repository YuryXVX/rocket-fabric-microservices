package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	Redis redisConfig `yaml:"redis"`
	Grpc  grpcConfig  `yaml:"grpc"`
	PG    pgConfig    `yaml:"pg"`
}

var appConfig *config

func MustLoad() {
	var cfg config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(fmt.Sprintf("не удалось загрузить конфиг: %v", err))
	}

	appConfig = &cfg
}

func AppConfig() *config {
	return appConfig
}
