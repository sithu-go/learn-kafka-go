package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	consume()
}

func consume() {
	config := sarama.NewConfig()
	config.ClientID = "myApp"
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	// config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	// config.Consumer.Group.Session.Timeout = 10 * time.Second   // Adjust as needed
	// config.Consumer.Group.Heartbeat.Interval = 3 * time.Second // Adjust as needed
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// brokers := []string{"kafka-0:9092", "kafka-1:9092", "kafka-2:9092"}
	brokers := []string{"localhost:9991", "localhost:9992", "localhost:9993"}
	groupID := "G1"
	topics := []string{"Users"}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Panicf("Error creating consumer group consumerGroup: %v", err)
	}

	defer consumerGroup.Close()

	fmt.Println("Connected")

	consumer := &Consumer{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, topics, consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				} else if errors.Is(err, sarama.ErrInconsistentGroupProtocol) {
					log.Fatalf("Error from consumer: %v", err)
				}

				log.Printf("Error Consuming : %v", err)
			}

			// if cancel will go here
			if ctx.Err() != nil {
				return
			}
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signals and close the consumer gracefully
	<-signals
	cancel()
	// select {
	// case sig := <-signals:
	// 	fmt.Printf("Caught signal %v; terminating\n", sig)
	// 	cancel()
	// }

	fmt.Println("Closing consumer")

}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
}

// 1st go in here
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error { return nil }

// 3rd go in here
// after all ConsumeClaim goroutines have exited but before the offsets are committed for the very last time.
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// 2nd go in here
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("Received message: %s at offset %d of partition %d\n", message.Value, message.Offset, message.Partition)
		session.MarkMessage(message, "")
	}
	return nil
}
