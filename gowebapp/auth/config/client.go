package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Client struct {
	BASE_URL string `default:"logfire.sh"`
	HOST     string `default:"localhost"`
	PORT     string `default:"3000"`
}

func CLIENT() Client {
	var client Client
	envconfig.MustProcess("Client", &client)
	log.Printf("Client variable: %v/n", &client)
	log.Println(&client)
	return client
}
