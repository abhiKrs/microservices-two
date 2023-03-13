package config

import (
	// log "web-api/app/utility/logger"

	"github.com/kelseyhightower/envconfig"
)

type Client struct {
	BASE_URL string `default:"beta.logfire.sh"`
	HOST     string `default:"localhost"`
	PORT     string `default:"3000"`
}

func CLIENT() Client {
	var client Client
	envconfig.MustProcess("Client", &client)
	// log.DebugLogger.Printf("Client variable: %v/n", &client)
	// log.Println(&client)
	return client
}
