package config

import (
	// "path/filepath"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Api struct {
	// ROOTPATH
	NAME              string
	HOST              string
	PORT              string
	ENV               string
	TLSHOST           string        `default:"beta.logfire.sh"`
	ReadHeaderTimeout time.Duration `default:"60s"`

	GracefulTimeout time.Duration `default:"8s"`

	REQUEST_LOG bool `default:"false"`

	MAILER_USERNAME    string
	MAILER_PASSWORD    string
	MAILER_SMTP_SERVER string
	MAILER_EMAIL       string `default:"support@logfire.sh"`

	AWS_ACCESS_KEY_ID     string `default:""`
	AWS_SECRET_ACCESS_KEY string `default:""`
}

func API() Api {
	var api Api
	// dir := filepath.Join("/public/emailTemplate").Dir("./main.go"). //Abs()
	// log.Println(dir)
	// log.Println("______________________________")
	envconfig.MustProcess("API", &api)
	log.Printf("Api variable: %v/n", &api)
	// log.Println(&api)

	return api
}
