package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

// NewKafkaProducer create the kafka producer connection
func NewKafkaProducer(log *zap.SugaredLogger, kafkaConfigMap *kafka.ConfigMap) *kafka.Producer {
	producer, err := kafka.NewProducer(kafkaConfigMap)
	if err != nil {
		log.Panicf("Error to create the kafka producer: %s", err)
	}

	log.Infof("Kafka Producer Connecteded: %v", producer.Len())
	return producer
}
