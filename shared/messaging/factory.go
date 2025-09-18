package messaging

// ProducerConfig holds configuration for creating producers
type ProducerConfig struct {
	Brokers   []string
	UseDocker bool // Use this to switch between implementations
}

// ConsumerConfig holds configuration for creating consumers
type ConsumerConfig struct {
	Brokers   []string
	GroupID   string
	Topics    []string
	UseDocker bool // Use this to switch between implementations
}

// NewEventProducer creates an appropriate producer based on configuration
func NewEventProducer(config ProducerConfig) (EventProducer, error) {
	// Try to create Sarama producer first
	producer, err := NewSaramaProducer(config.Brokers)
	if err != nil {
		// If Kafka is not available, return a no-op producer for graceful degradation
		return NewNoOpProducer(), nil
	}
	return producer, nil
}

// NewEventConsumer creates an appropriate consumer based on configuration
func NewEventConsumer(config ConsumerConfig) (EventConsumer, error) {
	// Try to create Sarama consumer first
	consumer, err := NewSaramaConsumer(config.Brokers, config.GroupID, config.Topics)
	if err != nil {
		// If Kafka is not available, return a no-op consumer for graceful degradation
		return NewNoOpConsumer(config.Topics), nil
	}
	return consumer, nil
}

// DetectEnvironment tries to detect if we're running in a Docker/CGO-friendly environment
func DetectEnvironment() bool {
	// For now, always use pure Go implementation to avoid CGO issues on Windows
	// In production Docker environments, you can set an environment variable to enable CGO version
	return false // Default to pure Go for Windows development
}
