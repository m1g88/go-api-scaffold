package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func Init() {
	if os.Getenv("ENV") == "prod" {
		return
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Process load environment variables into v.
func Process(v interface{}) error {
	return envconfig.Process("", v)
}

// MustProcess just wrap Process but it panic if Process has an error.
func MustProcess(prefix string, v interface{}) {
	envconfig.MustProcess(prefix, v)
}
