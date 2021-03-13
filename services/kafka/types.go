package kafka

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}
