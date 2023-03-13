package config

import (
	// log "web-api/app/utility/logger"

	"github.com/kelseyhightower/envconfig"
)

type Google struct {
	CLIENT_ID     string
	CLIENT_SECRET string
}

func GOOGLE() Google {
	var google Google
	envconfig.MustProcess("GOOGLE", &google)
	// log.DebugLogger.Printf("Google variable: %v/n", &google)
	return google
}
