package config

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type PGDatabase struct {
	DRIVER   string
	HOST     string
	PORT     uint16
	DB       string
	USER     string
	PASSWORD string
	SslMode  string `default:"disable"`
	// MaxConnectionPool      int           `default:"4"`
	// MaxIdleConnections     int           `default:"4"`
	ConnectionsMaxLifeTime time.Duration `default:"300s"`
}

func PGDataStore() PGDatabase {
	var db PGDatabase
	envconfig.MustProcess("POSTGRES", &db)
	log.Printf("Postgres variable: %v/n", &db)

	return db
}
