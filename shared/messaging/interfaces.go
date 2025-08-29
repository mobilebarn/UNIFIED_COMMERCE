package messaging

// EventProducer defines the interface for publishing messages
type EventProducer interface {
	Publish(topic, key string, payload []byte) error
	Close() error
}

// EventConsumer defines the interface for consuming messages
type EventConsumer interface {
	Subscribe(topics []string) error
	ReadMessage() (*Message, error)
	CommitMessage(msg *Message) error
	Close() error
}

// Message represents a message consumed from the event stream
type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Partition int32
	Offset    int64
}

// MessageHandler defines the interface for handling consumed messages
type MessageHandler interface {
	HandleMessage(msg *Message) error
}
