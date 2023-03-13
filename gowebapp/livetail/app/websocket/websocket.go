package websocket

import (
	"net/http"

	// log "livetail/app/utility/logger"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		return ws, err
// 	}
// 	return ws, nil
// }
