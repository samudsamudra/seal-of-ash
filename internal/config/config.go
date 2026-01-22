package config

import (
	"os"
	"github.com/joho/godotenv"
)

func Load() {
	_ = godotenv.Load()
}

func Get(key string) string {
	return os.Getenv(key)
}
