package config

import (
	log "web-api/app/utility/logger"

	"github.com/joho/godotenv"
)

type Config struct {
	Api        Api
	PGDatabase PGDatabase
	Client     Client
	Google     Google
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.ErrorLogger.Println(err)
	}

	// for _, e := range os.Environ() {

	// 	pair := strings.SplitN(e, "=", 2)
	// 	log.Printf("%s: %s\n", pair[0], pair[1])
	// }

	return &Config{
		Api:        API(),
		PGDatabase: PGDataStore(),
		Client:     CLIENT(),
		Google:     GOOGLE(),
	}
}

var CFG = New()
