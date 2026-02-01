package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, proceeding with OS environment variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}
	return &Config{
		Port: port,
	}
}
