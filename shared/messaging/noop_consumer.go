package messaging

import (
	"log"
)

// NoOpConsumer is a no-operation event consumer that logs instead of consuming from Kafka
// This allows services to function gracefully when Kafka is not available
type NoOpConsumer struct {
	topics []string
}

// NewNoOpConsumer creates a new no-op consumer
func NewNoOpConsumer(topics []string) EventConsumer {
	log.Println("Warning: Using no-op event consumer - events will be logged but not consumed from Kafka")
	return &NoOpConsumer{topics: topics}
}

// Subscribe logs the subscription instead of subscribing to Kafka
func (c *NoOpConsumer) Subscribe(topics []string) error {
	c.topics = topics
	log.Printf("No-op consumer: would subscribe to topics: %v", topics)
	return nil
}

// ReadMessage simulates reading messages by blocking forever
// In a real scenario, this would prevent the consumer loop from consuming anything
func (c *NoOpConsumer) ReadMessage() (*Message, error) {
	// Block forever to simulate waiting for messages that never come
	// This prevents the consumer loop from spinning and consuming CPU
	select {}
}

// CommitMessage does nothing for the no-op consumer
func (c *NoOpConsumer) CommitMessage(msg *Message) error {
	log.Printf("No-op consumer: would commit message from topic=%s", msg.Topic)
	return nil
}

// Close does nothing for the no-op consumer
func (c *NoOpConsumer) Close() error {
	return nil
}