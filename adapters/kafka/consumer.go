package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

// NewKafkaConsumer create the kafka consumer connection
func NewKafkaConsumer(log *zap.SugaredLogger, kafkaConfigMap *kafka.ConfigMap) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(kafkaConfigMap)
	if err != nil {
		log.Panicf("Error to create the kafka producer: %s", err)
	}

	log.Infof("Kafka Consumer Connecteded")
	return consumer
}

// CloseConsumer Close the kafka consumer
func CloseConsumer(consumer *kafka.Consumer) {
	_ = consumer.Close()
}
