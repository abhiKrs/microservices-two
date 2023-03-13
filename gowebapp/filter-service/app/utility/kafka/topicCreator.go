package kafka

import (
	log "filter-service/app/utility/logger"

	"github.com/Shopify/sarama"
)

func CreateTopic(topic string) error {
	log.DebugLogger.Println("started topic creation")
	brokers := []string{
		BootstrapServersDefault,
		// "kafka1:9092",
		// "kafka2:9092",
	}
	detail := sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_4_0_0
	admin, err := sarama.NewClusterAdmin(brokers, sarama.NewConfig())
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	err = admin.CreateTopic(topic, &detail, false)
	if err != nil {
		if err == sarama.ErrTopicAlreadyExists {
			log.InfoLogger.Println(err)
			return nil
		} else {
			log.ErrorLogger.Println("Error from kafka topic admin : ", err)
			return err
		}

	}
	return nil
}

func DeleteTopic(topic string) error {
	log.DebugLogger.Println("started topic deletion")
	brokers := []string{
		BootstrapServersDefault,
	}

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_4_0_0
	admin, err := sarama.NewClusterAdmin(brokers, sarama.NewConfig())
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	err = admin.DeleteTopic(topic)
	if err != nil {
		if err == sarama.ErrTopicDeletionDisabled {
			log.InfoLogger.Println(err)
			return nil
		} else {
			log.ErrorLogger.Println("Error from kafka topic admin : ", err)
			return err
		}

	}
	return nil
}
