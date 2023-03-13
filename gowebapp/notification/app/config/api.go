package config

import (
	// "path/filepath"

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

	MAILGUN_KEY    string
	MAILGUN_DOMAIN string

	AWS_ACCESS_KEY_ID     string `default:""`
	AWS_SECRET_ACCESS_KEY string `default:""`
}

func API() Api {
	var api Api
	envconfig.MustProcess("API", &api)
	return api
}
