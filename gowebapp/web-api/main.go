package main

import (
	"web-api/app/server"
	log "web-api/app/utility/logger"
)

var Version = "v0.1.0"

func main() {

	log.DebugLogger.Println("----hosting-----")
	log.DebugLogger.Printf("Starting API version: %s\n", Version)
	s := server.New()
	s.Init(Version)

	s.Run()
}
