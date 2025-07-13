package config

import (
	"log"
	"os"
)

var (
	ProcessorDefaultURL      string
	ProcessorFallbackURL     string
	PostgresDatabaseHostname string
	PostgresDatabasePort     string
	PostgresDatabaseName     string
	PostgresDatabaseUser     string
	PostgresDatabasePassword string
	RedisURL                 string
)

func Init() {
	ProcessorDefaultURL = os.Getenv("PROCESSOR_DEFAULT_URL")
	if ProcessorDefaultURL == "" {
		log.Fatal("PROCESSOR_DEFAULT_URL is not set")
	}

	ProcessorFallbackURL = os.Getenv("PROCESSOR_FALLBACK_URL")
	if ProcessorFallbackURL == "" {
		log.Fatal("PROCESSOR_FALLBACK_URL is not set")
	}

	PostgresDatabaseHostname = os.Getenv("POSTGRES_DATABASE_HOSTNAME")
	if PostgresDatabaseHostname == "" {
		log.Fatal("POSTGRES_DATABASE_HOSTNAME is not set")
	}

	PostgresDatabasePort = os.Getenv("POSTGRES_DATABASE_PORT")
	if PostgresDatabasePort == "" {
		log.Fatal("POSTGRES_DATABASE_PORT is not set")
	}

	PostgresDatabaseName = os.Getenv("POSTGRES_DATABASE_NAME")
	if PostgresDatabaseName == "" {
		log.Fatal("POSTGRES_DATABASE_NAME is not set")
	}

	PostgresDatabaseUser = os.Getenv("POSTGRES_DATABASE_USER")
	if PostgresDatabaseUser == "" {
		log.Fatal("POSTGRES_DATABASE_USER is not set")
	}

	PostgresDatabasePassword = os.Getenv("POSTGRES_DATABASE_PASSWORD")
	if PostgresDatabasePassword == "" {
		log.Fatal("POSTGRES_DATABASE_PASSWORD is not set")
	}

	RedisURL = os.Getenv("REDIS_URL")
	if RedisURL == "" {
		log.Fatal("REDIS_URL is not set")
	}

	log.Println("Configuration initialized successfully")
}
