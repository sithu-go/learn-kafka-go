package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	config := sarama.NewConfig()

	config.ClientID = "myApp"
	config.Version = sarama.V3_5_1_0

	brokers := []string{"localhost:9991", "localhost:9992", "localhost:9993"}

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka admin: %v", err)
	}
	defer admin.Close()

	fmt.Println("Connected...")

	// Use a signal channel to gracefully handle SIGINT
	// sigchan := make(chan os.Signal, 1)
	// signal.Notify(sigchan, os.Interrupt)

	// // Deleting Topic
	// if err := admin.DeleteTopic("Users"); err != nil {
	// 	log.Printf("Error deleting Topic: %v\n", err)
	// }

	// Create Topic
	// topicDetail := &sarama.TopicDetail{
	// 	NumPartitions:     2,
	// 	ReplicationFactor: 3,
	// }

	// // validateOnly
	// // determines whether the topic will actually be created or not.
	// //  If validateOnly is set to true, the request will be validated, but the topic will not be created.
	// if err := admin.CreateTopic("Users", topicDetail, false); err != nil {
	// 	log.Printf("Error creating topics: %v", err)
	// }

	// Listing Topics
	topics, err := admin.ListTopics()
	if err != nil {
		log.Printf("Error listing topics: %v \n", err)
	}
	fmt.Printf("%+v\n", topics)

	// Listing Groups
	groups, err := admin.ListConsumerGroups()
	if err != nil {
		log.Printf("Error listing consumer groups: %v \n", err)
	}

	fmt.Printf("%+v\n", groups)

	// topicStatus, err := admin.ListPartitionReassignments("Users", []int32{0, 1})
	// if err != nil {
	// 	log.Printf("Error Listing PartitionReassignments: %v \n", err)
	// }

	// fmt.Printf("%+v\n", topicStatus)

}
