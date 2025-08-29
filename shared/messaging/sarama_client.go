package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

// SaramaProducer implements EventProducer using Shopify's sarama library (pure Go)
type SaramaProducer struct {
	producer sarama.SyncProducer
}

// NewSaramaProducer creates a new pure Go Kafka producer
func NewSaramaProducer(brokers []string) (*SaramaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama producer: %w", err)
	}

	return &SaramaProducer{producer: producer}, nil
}

// Publish sends a message to a Kafka topic
func (p *SaramaProducer) Publish(topic, key string, payload []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(payload),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Printf("Message sent to topic %s, partition %d, offset %d\n", topic, partition, offset)
	return nil
}

// Close closes the producer
func (p *SaramaProducer) Close() error {
	return p.producer.Close()
}

// SaramaConsumer implements EventConsumer using Shopify's sarama library (pure Go)
type SaramaConsumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	groupID  string
	messages chan *Message
	handler  *ConsumerGroupHandler
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewSaramaConsumer creates a new pure Go Kafka consumer
func NewSaramaConsumer(brokers []string, groupID string, topics []string) (*SaramaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama consumer: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	messages := make(chan *Message, 100)
	handler := &ConsumerGroupHandler{messages: messages}

	return &SaramaConsumer{
		consumer: consumer,
		topics:   topics,
		groupID:  groupID,
		messages: messages,
		handler:  handler,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Subscribe subscribes to the configured topics and starts consumption
func (c *SaramaConsumer) Subscribe(topics []string) error {
	c.topics = topics
	
	// Start consuming in a goroutine
	go func() {
		for {
			if err := c.consumer.Consume(c.ctx, c.topics, c.handler); err != nil {
				if err == sarama.ErrClosedConsumerGroup {
					return
				}
				fmt.Printf("Error from consumer: %v\n", err)
			}
			
			// Check if context was cancelled
			if c.ctx.Err() != nil {
				return
			}
		}
	}()
	
	return nil
}

// ReadMessage reads a message from the consumer
func (c *SaramaConsumer) ReadMessage() (*Message, error) {
	select {
	case msg := <-c.messages:
		return msg, nil
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	}
}

// CommitMessage commits a message offset
func (c *SaramaConsumer) CommitMessage(msg *Message) error {
	// Message commit is handled automatically in the ConsumerGroupHandler
	return nil
}

// Close closes the consumer
func (c *SaramaConsumer) Close() error {
	c.cancel()
	return c.consumer.Close()
}

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct {
	messages chan *Message
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		msg := &Message{
			Topic:     message.Topic,
			Key:       message.Key,
			Value:     message.Value,
			Partition: message.Partition,
			Offset:    message.Offset,
		}
		
		select {
		case h.messages <- msg:
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
	return nil
}
