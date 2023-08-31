package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	config.Producer.Partitioner = sarama.NewManualPartitioner

	brokers := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	defer producer.Close()

	fmt.Println("Producing data ...")

	// A-M 0, N-Z 1
	// For string
	partition := int32(0)

	if string(strings.ToUpper(msg)[0]) >= "N" {
		partition = int32(1)
		fmt.Println(partition)
	}

	message := &sarama.ProducerMessage{
		Topic:     "Users",
		Partition: partition,
		Value:     sarama.StringEncoder(msg),
	}

	// for number
	// numberToSend := 42 // Replace with the number you want to send

	// var buf bytes.Buffer
	// if err := binary.Write(&buf, binary.BigEndian, int32(numberToSend)); err != nil {
	// 	log.Fatalf("Error encoding number: %v", err)
	// }

	// message := &sarama.ProducerMessage{
	// 	Topic: "Users",
	// 	Value: sarama.ByteEncoder(buf.Bytes()), // Use ByteEncoder to send the encoded byte slice
	// }

	// for JSON
	// jsonData := map[string]interface{}{
	// 	"key1": "value1",
	// 	"key2": 42,
	// }

	// encodedData, err := json.Marshal(jsonData)
	// if err != nil {
	// 	log.Fatalf("Error encoding JSON: %v", err)
	// }

	// message := &sarama.ProducerMessage{
	// 	Topic: "Users",
	// 	Value: sarama.ByteEncoder(encodedData),
	// }

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Error sending Message: %v", err)
	}

	fmt.Printf("parttion = %d, offset = %d \n", partition, offset)

}
