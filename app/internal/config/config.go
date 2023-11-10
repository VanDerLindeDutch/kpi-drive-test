package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Kpi struct {
		BaseUrl  string `env:"KPI_BASEURL" env-default:"https://development.kpi-drive.ru/_api/"`
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
