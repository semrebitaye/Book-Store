package initializers

import (
	"log"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		slog.Error("Error loading .env file", err)
	}
}
