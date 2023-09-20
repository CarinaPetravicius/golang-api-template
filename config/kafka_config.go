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
	plain         = "PLAIN"
)

// NewKafkaConfigMap config the kafka connection properties
func NewKafkaConfigMap(log *zap.SugaredLogger, config KafkaConfiguration) *kafka.ConfigMap {
	var kafkaConf = &kafka.ConfigMap{
		"bootstrap.servers": config.Servers,
		"message.max.bytes": 1000000,
		"retries":           5,
		"retry.backoff.ms":  1000,
		// Acks property controls how many partition replicas must acknowledge the receipt of a record before a producer can consider a particular write operation as successful.
		// acks = -1, the producer waits for the ack. Having the messages replicated to all the partition replicas.
		// acks = 1, the leader must receive the record and respond before the write is considered successful
		// acks = 0, the write is considered successful the moment the request is sent out. No need to wait for a response.
		"acks": -1,
	}

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

	return kafkaConf
}

func setSSLProperties(kafkaConf *kafka.ConfigMap, config *KafkaConfiguration) {
	_ = kafkaConf.SetKey("sasl.mechanism", plain)
	_ = kafkaConf.SetKey("sasl.username", config.User)
	_ = kafkaConf.SetKey("sasl.password", config.Pass)
}
