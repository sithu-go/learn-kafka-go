package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	config := sarama.NewConfig()

	config.ClientID = "myApp"

	brokers := []string{"localhost:9092"}

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka admin: %v", err)
	}
	defer admin.Close()

	fmt.Println("Connected...")

	// // Deleting Topic
	// if err := admin.DeleteTopic("Users"); err != nil {
	// 	log.Printf("Error deleting Topic: %v\n", err)
	// }

	// Create Topic
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     2,
		ReplicationFactor: -1,
	}

	if err := admin.CreateTopic("Users", topicDetail, false); err != nil {
		log.Fatalf("Error creating topics: %v", err)
	}

}
