package kafka

import "github.com/knoguchi/go_project_template/services"

type Config struct {
	services.ServiceConfig
	Brokers  []string       `json:"brokers"`
	TLS      bool           `json:"tls"`
	Consumer ConsumerConfig `json:"consumer"`
}

type ConsumerConfig struct {
	Topics []string `json:"topics"`
	Group  string   `json:"group"`
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}
