package config

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const (
	protocol      = "security.protocol"
	plaintext     = "plaintext"
	saslSsl       = "sasl_ssl"
	saslPlaintext = "sasl_plaintext"
)

// NewKafkaProducer init the kafka producer connection
func NewKafkaProducer(log *zap.SugaredLogger, config KafkaConfiguration) *kafka.Producer {
	var kafkaConf = &kafka.ConfigMap{
		"message.max.bytes": 1000000,
		"retries":           5,
		"retry.backoff.ms":  1000,
		// Acks property controls how many partition replicas must acknowledge the receipt of a record before a producer can consider a particular write operation as successful.
		// acks = -1, the producer waits for the ack. Having the messages replicated to all the partition replicas.
		// acks = 1, the leader must receive the record and respond before the write is considered successful
		// acks = 0, the write is considered successful the moment the request is sent out. No need to wait for a response.
		"acks": -1,
	}

	_ = kafkaConf.SetKey("bootstrap.servers", config.Servers)

	switch config.SecurityProtocol {
	case plaintext:
		_ = kafkaConf.SetKey(protocol, plaintext)
	case saslSsl:
		_ = kafkaConf.SetKey(protocol, saslSsl)
		_ = kafkaConf.SetKey("ssl.ca.location", "conf/ca-cert.pem")
		setSSLProperties(kafkaConf, &config)
	case saslPlaintext:
		_ = kafkaConf.SetKey(protocol, saslPlaintext)
		setSSLProperties(kafkaConf, &config)
	default:
		log.Panic(kafka.NewError(kafka.ErrUnknownProtocol, "Unknown kafka protocol", true))
	}

	producer, err := kafka.NewProducer(kafkaConf)
	if err != nil {
		log.Panicf("Error to create the kafka producer: %s", err)
	}

	log.Infof("Kafka Producer Connecteded: %v", producer.Len())
	return producer
}

func setSSLProperties(kafkaConf *kafka.ConfigMap, config *KafkaConfiguration) {
	_ = kafkaConf.SetKey("sasl.mechanism", "PLAIN")
	_ = kafkaConf.SetKey("sasl.username", config.User)
	_ = kafkaConf.SetKey("sasl.password", config.Pass)
}
