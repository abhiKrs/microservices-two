package config

import (
	"time"

	// log "web-api/app/utility/logger"
	"github.com/kelseyhightower/envconfig"
)

type Api struct {
	// ROOTPATH
	NAME              string
	HOST              string
	PORT              string
	ENV               string
	TLSHOST           string        `default:"apibeta.logfire.sh"`
	ReadHeaderTimeout time.Duration `default:"60s"`

	GracefulTimeout time.Duration `default:"8s"`

	REQUEST_LOG bool `default:"false"`
}

func API() Api {
	var api Api
	// dir := filepath.Join("/public/emailTemplate").Dir("./main.go"). //Abs()
	// log.Println(dir)
	// log.Println("______________________________")
	envconfig.MustProcess("API", &api)
	// log.DebugLogger.Printf("Api variable: %v/n", &api)
	// log.Println(&api)

	return api
}
