package kafka

// IMessage kafka message interface
type IMessage interface {
	ProduceMessage(topicName, value, eventName, traceID string)
}
