package main

import (
	// "bufio"
	// "encoding/json"
	"fmt"
	"net/http"

	// "fmt"
	// "os"
	// "os/signal"
	// "syscall"
	// "time"

	"stream_source/app/api/controller"
	// "stream_source/app/kafkaproducer"
	"stream_source/app/logs/model"
	// "stream_source/app/logs/utility"
	log "stream_source/app/utility/logger"
	// "github.com/Shopify/sarama"
)

type Message struct {
	Data model.Data
}

func main() {
	port := ":4000"
	manager := controller.NewLogSourceController()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
		log.InfoLogger.Println("Started Server")
	})
	http.HandleFunc("/logfire.sh", manager.SendLog)
	http.ListenAndServe(port, nil)
	log.InfoLogger.Printf("Server listening at port: %v", port)
}

// func main() {
// 	config := kafkaproducer.NewProducerConfig()
// 	log.InfoLogger.Printf("Go producer starting with config=%+v\n", config)

// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

// 	producerConfig := sarama.NewConfig()
// 	producerConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.ProducerAcks)
// 	producerConfig.Producer.Return.Successes = true
// 	producer, err := sarama.NewSyncProducer([]string{config.BootstrapServers}, producerConfig)
// 	if err != nil {
// 		log.ErrorLogger.Printf("Error creating the Sarama sync producer: %v", err)
// 		os.Exit(1)
// 	}

// 	end := make(chan int, 1)
// 	go func() {
// 		for {
// 			log.InfoLogger.Println("streaming")
// 			// file, ferr := os.Open("dummy_logs.txt")
// 			file, ferr := os.Open("logs_message.txt")
// 			if ferr != nil {
// 				log.ErrorLogger.Panic(ferr)
// 			}
// 			scanner := bufio.NewScanner(file)
// 			for scanner.Scan() {
// 				line := scanner.Text()
// 				data := utility.GenerateData(line)
// 				json, err := json.Marshal(data)
// 				if err != nil {
// 					log.ErrorLogger.Panic(ferr)
// 					break
// 				}
// 				msg := &sarama.ProducerMessage{
// 					Topic: config.Topic,
// 					Value: sarama.ByteEncoder(json),
// 				}
// 				partition, offset, err := producer.SendMessage(msg)
// 				if err != nil {
// 					log.ErrorLogger.Printf("Erros sending message: %v\n", err)
// 				} else {
// 					log.DebugLogger.Printf("Message sent: partition=%d, offset=%d\n", partition, offset)
// 				}
// 				time.Sleep(time.Second)
// 			}

// 			log.InfoLogger.Printf("Message Received: %+v\n", "logs_message.txt")

// 		}
// 	}()

// 	// waiting for the end of all messages sent or an OS signal
// 	select {
// 	case <-end:
// 		log.DebugLogger.Printf("Finished to send %d messages\n", config.MessageCount)
// 	case sig := <-signals:
// 		log.DebugLogger.Printf("Got signal: %v\n", sig)
// 	}

// 	err = producer.Close()
// 	if err != nil {
// 		log.ErrorLogger.Printf("Error closing the Sarama sync producer: %v", err)
// 		os.Exit(1)
// 	}
// 	log.InfoLogger.Printf("Producer closed")
// }
