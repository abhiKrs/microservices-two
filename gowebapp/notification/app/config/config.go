package config

import (
	log "notification/app/utility"
	// "os"
	// "strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Api    Api
	Client Client
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.ErrorLogger.Println(err)
	}

	return &Config{
		Api:    API(),
		Client: CLIENT(),
	}
}

var CFG = New()
