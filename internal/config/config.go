package config

import (
	"log"
	"os"
)

var (
	ProcessorDefaultURL  string
	ProcessorFallbackURL string
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
}

