package controller

import (
	"encoding/json"
	"net/http"
	"os"

	"stream_source/app/kafkaproducer"
	log "stream_source/app/utility/logger"

	"github.com/Shopify/sarama"
)

type Message struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

type LogSourceController struct {
	// sourceService service.SourceService
	// validator     *validator.Validate
}

func NewLogSourceController() *LogSourceController {
	return &LogSourceController{
		// sourceService: *sourceService,
		// validator:     validator,
	}
}

func (pc *LogSourceController) SendLog(w http.ResponseWriter, r *http.Request) {

	// reqToken := r.Header.Get("Authorization")
	// if reqToken == "" {
	// 	log.ErrorLogger.Println("Missing Authorization Header")
	// 	respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	// }

	// TODO use JWT token
	// splitToken := strings.Split(reqToken, "Bearer ")
	// sourceId := splitToken[1]
	// log.InfoLogger.Println("sourceId")

	// var reqBody models.SourceLog
	var reqBody interface{}
	// err := reqBody.Bind(r.Body)
	// if err != nil {
	// 	log.ErrorLogger.Println(err)
	// 	respond.Error(w, http.StatusBadRequest, err)
	// 	return
	// }
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.WarningLogger.Println(err)
		// http.Error(w, err.Error(), http.StatusBadRequest)
		// return
	}

	log.DebugLogger.Println(reqBody)

	// errs := validate.Validate(pc.validator, reqBody)
	// if errs != nil {
	// 	respond.Errors(w, http.StatusBadRequest, errs)
	// 	return
	// }

	// if sourceId == "" {
	// 	response := models.SourceLogResponse{IsSuccessful: false, Message: []string{"Wrong or missing profile Id"}}
	// 	respond.Json(w, http.StatusBadRequest, response)
	// 	return
	// }

	config := kafkaproducer.NewProducerConfig()
	log.InfoLogger.Printf("Go producer starting with config=%+v\n", config)

	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.ProducerAcks)
	producerConfig.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{config.BootstrapServers}, producerConfig)
	if err != nil {
		log.ErrorLogger.Printf("Error creating the Sarama sync producer: %v", err)
		os.Exit(1)
	}

	end := make(chan int, 1)
	go func() {
		// for i := int64(0); i < config.MessageCount; i++ {
		// for i, logData := range reqBody.Data {
		// log.DebugLogger.Printf("%v-%d", logData, i)
		message := Message{Type: 3, Data: reqBody}
		byteMsg, err := json.Marshal(message)
		if err != nil {
			log.ErrorLogger.Println(err)
		}

		msg := &sarama.ProducerMessage{
			Topic: "my-newtopic",
			Value: sarama.ByteEncoder(byteMsg),
		}
		log.DebugLogger.Printf("Sending message: value=%s\n", msg.Value)
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.DebugLogger.Printf("Erros sending message: %v\n", err)
		} else {
			log.DebugLogger.Printf("Message sent: partition=%d, offset=%d\n", partition, offset)
		}
		end <- 1
	}()
}
