package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

var wg sync.WaitGroup

func main() {
	consume()
}

func consume() {
	config := sarama.NewConfig()
	config.ClientID = "myapp"
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	kafkaBrokers := []string{"localhost:9092"}
	consumer, err := sarama.NewConsumer(kafkaBrokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	fmt.Println("Connected")

	topic := "Users"
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Error fetching partitions: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background()) // Create a context

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			log.Fatalf("Error creating partition consumer: %v", err)
		}
		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			defer partitionConsumer.Close()

			for {
				select {
				case <-ctx.Done(): // Check if context is done (signal received)
					fmt.Println("Stopping consumer...")
					return
				case msg := <-pc.Messages():
					fmt.Printf("Received Message: %s on Partition %d at Offset %d\n",
						msg.Value, msg.Partition, msg.Offset)
				}
			}
		}(partitionConsumer)
	}

	<-sigchan // Wait for a signal to exit

	fmt.Println("Stopping all consumers...")
	cancel()  // Cancel the context to signal all consumers to stop
	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("All consumers stopped.")
}
