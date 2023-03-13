package kafka

import (
	"encoding/json"

	logs "web-api/app/utility/logger"

	// "github.com/Shopify/sarama"
	// "github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	// "github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	// "github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	// "github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

func SendLog(kafkaTopic string, payload interface{}) error {
	// Produce a new record to the topic...
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": BootstrapServersDefault,
		// "sasl.mechanisms":   "PLAIN",
		// "security.protocol": "SASL_PLAINTEXT",
		// "sasl.username":     ccloudAPIKey,
		// "sasl.password":     ccloudAPISecret
	})

	if err != nil {
		logs.ErrorLogger.Printf("Failed to create producer: %s", err)
		return err
	}

	bytePayload, err := json.Marshal(payload)
	if err != nil {
		logs.ErrorLogger.Println(err)
		return err
	}

	producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &kafkaTopic,
				Partition: 0,
			},
			Value: bytePayload,
		},
		nil,
	)

	// Wait for delivery report
	e := <-producer.Events()

	message := e.(*kafka.Message)
	if message.TopicPartition.Error != nil {
		logs.WarningLogger.Printf("failed to deliver message: %v\n", message.TopicPartition)
	} else {
		logs.WarningLogger.Printf("delivered to topic %s [%d] at offset %v\n",
			*message.TopicPartition.Topic,
			message.TopicPartition.Partition,
			message.TopicPartition.Offset)
	}

	producer.Close()
	return nil
}

// func SendLog(kafkaTopic string, msg interface{}) error {
// 	config := NewProducerConfig()
// 	logs.InfoLogger.Printf("Go producer starting with config=%+v\n", config)

// 	producerConfig := sarama.NewConfig()
// 	producerConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.ProducerAcks)
// 	producerConfig.Producer.Return.Successes = true
// 	producer, err := sarama.NewSyncProducer([]string{config.BootstrapServers}, producerConfig)
// 	if err != nil {
// 		logs.ErrorLogger.Printf("Error creating the Sarama sync producer: %v", err)
// 		// os.Exit(1)
// 		return err
// 	}

// 	go func() {
// 		// for i, logData := range msg.Data {
// 		logs.DebugLogger.Println(msg)
// 		// message := Message{Type: 3, Data: msg}
// 		// byteMsg, err := json.Marshal(message)
// 		byteMsg, err := json.Marshal(msg)
// 		if err != nil {
// 			logs.ErrorLogger.Println(err)

// 			return
// 		}

// 		msg := &sarama.ProducerMessage{
// 			Topic: kafkaTopic,
// 			Value: sarama.ByteEncoder(byteMsg),
// 		}
// 		logs.DebugLogger.Printf("Sending message: value=%s\n", msg.Value)
// 		partition, offset, err := producer.SendMessage(msg)
// 		if err != nil {
// 			logs.DebugLogger.Printf("Erros sending message: %v\n", err)
// 			return
// 		} else {
// 			logs.DebugLogger.Printf("Message sent: partition=%d, offset=%d\n", partition, offset)
// 			return
// 		}
// 	}()

// 	return nil
// }
