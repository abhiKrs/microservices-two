package main

import (
	"net/http"

	log "livetail/app/utility/logger"
	"livetail/app/websocket"
)

func setupRoutes() {
	// manager := websocket.NewManager()
	// http.HandleFunc("/ws", manager.ServeWs)
	// http.HandleFunc("/ws", manager.ServeWsKafka)
	http.HandleFunc("/ws", websocket.ServeClient)
}

func main() {
	port := ":4000"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello from /livetail"))
	})

	setupRoutes()

	log.InfoLogger.Printf("Server listening at port: %v", port)

	http.ListenAndServe(port, nil)

}
