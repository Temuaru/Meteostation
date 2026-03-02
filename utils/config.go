package utils

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Определение пути к файлу переменных окружения .env
const envPath = "./.env"

// Определение структуры конфигурации
type Config struct {
	WebPort       string `env:"WEB_PORT"`
	DBUser        string `env:"DB_USER"`
	DBPassword    string `env:"DB_PASSWORD"`
	DBName        string `env:"DB_NAME"`
	DBHost        string `env:"DB_HOST"`
	DBSSLMode     string `env:"DB_SSLMODE"`
	DBPort        int    `env:"DB_PORT"`
	CollectorHost string `env:"COLLECTOR_HOST"`
}

// Заполняем все необходимые поля конфигурации
func (c *Config) LoadConfig() {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("error reading %s: %s", envPath, err)
	}

	err = env.Parse(c)
	if err != nil {
		log.Printf("Error parsing env to struct: %s", err)
	}

}
