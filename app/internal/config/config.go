package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME" env-default:"bis"`
	IsDebug     bool   `env:"IS_DEBUG" env-default:"true"`
	Listen      struct {
		BindIP string `env:"BIND_IP" env-default:"0.0.0.0"`
		Port   string `env:"PORT" env-default:"8081"`
	}
	DatabaseConfig struct {
		Port     int    `env:"DB_PORT" env-default:"5435"`
		Host     string `env:"DB_HOST" env-default:"127.0.0.1"`
		Name     string `env:"DB_NAME" env-default:"bis"`
		User     string `env:"DB_USER" env-default:"postgres"`
		Password string `env:"DB_PASSWORD" env-default:"postgres"`
	}
	Timezone struct {
		Local string `env:"TIMEZONE_LOCAL" env-default:"Europe/Riga"`
	}
	JWT struct {
		CertificatePath string `env:"JWT_CERTIFICATE_PATH" env-default:"var/jwt-certificate.pem"`
	}
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
