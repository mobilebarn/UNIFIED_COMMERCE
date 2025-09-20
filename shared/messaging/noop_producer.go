package messaging

import (
	"log"
)

// NoOpProducer is a no-operation event producer that logs messages instead of sending them
// This allows services to function gracefully when Kafka is not available
type NoOpProducer struct{}

// NewNoOpProducer creates a new no-op producer
func NewNoOpProducer() EventProducer {
	log.Println("Warning: Using no-op event producer - events will be logged but not sent to Kafka")
	return &NoOpProducer{}
}

// Publish logs the message instead of sending to Kafka
func (p *NoOpProducer) Publish(topic, key string, payload []byte) error {
	log.Printf("No-op producer: would publish to topic=%s, key=%s, payload_size=%d bytes",
		topic, key, len(payload))
	return nil
}

// Close does nothing for the no-op producer
func (p *NoOpProducer) Close() error {
	return nil
}
