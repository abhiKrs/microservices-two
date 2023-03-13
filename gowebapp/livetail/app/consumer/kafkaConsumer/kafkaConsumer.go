package kafkaConsumer

import (

	// "os/signal"
	// "strings"
	logs "livetail/app/utility/logger"

	"github.com/Shopify/sarama"
)

var (
	// Brokers = []string{"my-cluster-kafka-bootstrap.kafka:9092"}
	Brokers = []string{"my-cluster-kafka-bootstrap:9092"}
)

func NewKafkaConfiguration(clientId string) *sarama.Config {
	conf := sarama.NewConfig()
	conf.ClientID = clientId
	conf.Consumer.Return.Errors = true
	return conf
}

func Consume(
	topic *string,
	master sarama.Consumer,
	doneChan chan bool,
	consumeChan chan *sarama.ConsumerMessage,
	errorChan chan *sarama.ConsumerError,
) {

	partitions, _ := master.Partitions(*topic)
	// this only consumes partition no 1, you would probably want to consume all partitions
	consumer, err := master.ConsumePartition(*topic, partitions[0], sarama.OffsetNewest)
	if err != nil {
		logs.ErrorLogger.Printf("Topic %v Partitions: %v", topic, partitions)
		logs.ErrorLogger.Println(err)
		// panic(err)
		return
	}

	logs.DebugLogger.Println(" Start consuming topic ", *topic)
	go func(consumer sarama.PartitionConsumer) {
		for {
			select {
			case consumerError := <-consumer.Errors():

				logs.WarningLogger.Println("consumerError: ", consumerError.Err)
				errorChan <- consumerError
				// logs.ErrorLogger.Println("consumerError: ", consumerError.Err)

			case msg := <-consumer.Messages():
				consumeChan <- msg
				// logs.ErrorLogger.Println("Got message on topic ", topic, string(msg.Value))
			case <-doneChan:
				// close(doneChan)
				logs.ErrorLogger.Println("Got Done Channel. Closing Kafka consumer channel")
				consumer.Close()
				return
			}
		}
	}(consumer)
}
