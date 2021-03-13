package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/knoguchi/go_project_template/services"
	"golang.org/x/sync/errgroup"
)

type KafkaService struct {
	services.Service
	// Brokers is Kafka bootstrap brokers to connect to, as a comma separated list
	Brokers   []string
	TLSEnable bool
	// Version is Kafka cluster version, default 2.1.1.
	Version string
	// Group is Kafka consumer group definition
	Group string
	// Assignor is consumer group partition assignment strategy (range, roundrobin, sticky)
	Assignor string
	// Oldest is Kafka consumer consume initial offset from oldest
	Oldest bool
	//
	Topics []string
}

func New() *KafkaService {
	svc := &KafkaService{
		Brokers:   []string{"localhost:9092"},
		TLSEnable: true,
		Topics:    []string{"test"},
		Group:     "dev",
		Version:   "2.1.1",
		Assignor:  "roundrobin",
	}

	return svc
}

func (k *KafkaService) Configure() {

}
func (k *KafkaService) Start(ctx context.Context) error {
	g, gctx := errgroup.WithContext(ctx)
	if k.Verbose {
		sarama.Logger = log
	}

	version, err := sarama.ParseKafkaVersion(k.Version)
	if err != nil {
		log.Errorf("Error parsing Kafka version: %v", err)
		return err
	}

	/**
	 * Setup a new Sarama consumer group
	 */
	config := sarama.NewConfig()
	config.Version = version

	if k.TLSEnable {
		config.Net.TLS.Enable = k.TLSEnable
	}

	switch k.Assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", k.Assignor)
	}

	if k.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	consumer := Consumer{
		ready: make(chan bool),
	}

	//ctx, cancel := context.WithCancel(gctx)
	client, err := sarama.NewConsumerGroup(k.Brokers, k.Group, config)
	if err != nil {
		log.Errorf("Error creating consumer group client: %v", err)
		return err
	}

	g.Go(func() error {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(gctx, k.Topics, &consumer); err != nil {
				log.Errorf("Error from consumer: %v", err)
				return gctx.Err()
			}
			// check if context was cancelled, signaling that the consumer should stop
			if gctx.Err() != nil {
				return gctx.Err()
			}
			consumer.ready = make(chan bool)
		}
	})

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	// wait for all errgroup goroutines
	err = g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Info("context was canceled")
		} else {
			log.Infof("received error: %v", err)
		}
	} else {
		log.Infoln("finished clean")
	}
	log.Info("kafka done")
	return nil
}

func (k *KafkaService) Status() error {
	return nil
}

func (k *KafkaService) Stop() error {
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}
