package websocket

import (
	// "bytes"
	"encoding/json"
	// "io"
	log "livetail/app/utility/logger"
	"net/http"
	// "livetail/app/utility/redis"
	// "github.com/Shopify/sarama"
)

func ServeClient(w http.ResponseWriter, r *http.Request) {
	log.InfoLogger.Println("New Kafka Connection")

	// var req *Filter
	// // err := json.NewDecoder(r.Body).Decode(r)
	// err := req.Bind(r.Body)
	// if err != nil {
	// 	log.DebugLogger.Println(err)
	// 	res := Response{IsSuccessful: false, Message: []string{err.Error()}}
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	byteRes, err := json.Marshal(res)
	// 	if err != nil {
	// 		// w.Write([]bytes([]string{err.Error()}))
	// 		return
	// 	}
	// 	w.Write(byteRes)
	// 	return
	// }

	// upgrade to websocket connection
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		res := Response{IsSuccessful: false, Message: []string{err.Error()}}
		byteRes, _ := json.Marshal(res)
		w.Write(byteRes)
		return
	}

	// create new client
	client := NewClient(conn, r)

	go client.ReadMessage()
	go client.WriteMessage()
	go client.StreamKafka()

}
