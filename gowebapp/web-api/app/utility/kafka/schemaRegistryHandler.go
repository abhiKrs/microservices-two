package kafka

import (
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io"
	// "net/http"
	// "strings"

	// "web-api/app/source/models"
	// "web-api/app/utility/kafka"
	logs "web-api/app/utility/logger"
	// "web-api/app/utility/respond"

	// "web-api/app/utility/validate"
	// "github.com/Shopify/sarama"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

type RecordsHandler struct {
	Producer *kafka.Producer
	// Ser      *avro.SpecificSerializer
	Ser   *avro.GenericSerializer
	Topic string
}

func NewRecordHandler(sourceTopic string) (*RecordsHandler, error) {
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": BootstrapServersDefault})

	if err != nil {
		logs.WarningLogger.Printf("Failed to create producer: %s\n", err)
		return nil, err
	}

	logs.DebugLogger.Printf("Created Producer %v\n", kafkaProducer)

	schemaClient, err := schemaregistry.NewClient(schemaregistry.NewConfig(SchemaRegistryUrl))

	if err != nil {
		logs.WarningLogger.Printf("Failed to create schema registry client: %s\n", err)
		return nil, err
	}

	// schemaClient.

	ser, err := avro.NewGenericSerializer(schemaClient, serde.ValueSerde, avro.NewSerializerConfig())
	// ser, err := avro.NewSpecificSerializer(schemaClient, serde.ValueSerde, avro.NewSerializerConfig())

	if err != nil {
		logs.WarningLogger.Printf("Failed to create serializer: %s\n", err)
		return nil, err
	}

	rh := &RecordsHandler{
		kafkaProducer,
		ser,
		sourceTopic,
	}
	return rh, nil
}
