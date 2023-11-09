package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Kpi struct {
		Username string `env:"KPI_USERNAME" env-default:"admin"`
		Password string `env:"KPI_PASSWORD" env-default:"admin"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadEnv(instance); err != nil {
			log.Fatal(err)
		}

	})
	return instance
}
