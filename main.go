package main

import (
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

func main() {
	msg := os.Args[1]
	produce(msg)
}

func produce(msg string) {
	config := sarama.NewConfig()
	config.ClientID = "myApp"
	config.Producer.Return.Successes = true

	brokers := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	defer producer.Close()

	fmt.Println("Producing data ...")

	// here's I set at 0
	partition := int32(0)

	message := &sarama.ProducerMessage{
		Topic:     "Users",
		Partition: partition,
		Value:     sarama.StringEncoder(msg),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Error sending Message: %v", err)
	}

	fmt.Printf("parttion = %d, offset = %d \n", partition, offset)

}
