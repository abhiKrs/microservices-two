package kafka

import (
	"encoding/json"

	log "filter-service/app/utility/logger"

	"github.com/Shopify/sarama"
)

func SendLog(kafkaTopic string, msg interface{}) error {
	config := NewProducerConfig()
	log.InfoLogger.Printf("Go producer starting with config=%+v\n", config)

	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.ProducerAcks)
	producerConfig.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{config.BootstrapServers}, producerConfig)
	if err != nil {
		log.ErrorLogger.Printf("Error creating the Sarama sync producer: %v", err)
		// os.Exit(1)
		return err
	}
	errs := make(chan error, 1)
	go func() {
		// for i, logData := range msg.Data {
		defer close(errs)
		log.DebugLogger.Println(msg)
		byteMsg, err := json.Marshal(msg)
		if err != nil {
			log.ErrorLogger.Println(err)
			errs <- err
			return
		}

		msg := &sarama.ProducerMessage{
			Topic: kafkaTopic,
			Value: sarama.ByteEncoder(byteMsg),
		}
		log.DebugLogger.Printf("Sending message: value=%s\n", msg.Value)
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.DebugLogger.Printf("Erros sending message: %v\n", err)
			errs <- err
			return
		} else {
			log.DebugLogger.Printf("Message sent: partition=%d, offset=%d\n", partition, offset)
			errs <- err
			return
		}

	}()

	// if

	return nil
}
