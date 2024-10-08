package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type EnvConfig struct {
	Port             string
	ConnectionString string
}

var AppConfig *EnvConfig

func LoadAppConfig(logger *zerolog.Logger) {
	logger.Info().Msg("Loading Server Configurations...")

	// comment this out when running in debug mode
	if err := godotenv.Load(); err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}

	envPort := os.Getenv("PORT")
	envDBString := os.Getenv("DB_STRING")

	AppConfig = &EnvConfig{
		Port:             envPort,
		ConnectionString: envDBString,
	}
}
