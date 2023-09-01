package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	// msg := os.Args[1]
	// produce(msg)
	dd()
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

// Bitwise test for
// log.SetFlags(log.LstdFlags | log.Lshortfile)
// func main() {
// 	SetFlags(20 | 21)

// }

// func SetFlags(flag int) {
// 	fmt.Println(flag)
// }

func dd() {
	// json
	b := []byte{123, 34, 102, 114, 111, 109, 34, 58, 123, 34, 105, 100, 34, 58, 49, 44, 34, 110, 97, 109, 101, 34, 58, 34, 69, 109, 109, 97, 34, 125, 44, 34, 116, 111, 34, 58, 123, 34, 105, 100, 34, 58, 50, 44, 34, 110, 97, 109, 101, 34, 58, 34, 66, 114, 117, 110, 111, 34, 125, 44, 34, 109, 101, 115, 115, 97, 103, 101, 34, 58, 34, 69, 109, 109, 97, 32, 115, 116, 97, 114, 116, 101, 100, 32, 102, 111, 108, 108, 111, 119, 105, 110, 103, 32, 121, 111, 117, 34, 125}

	// number uint32
	// i := []byte{0, 0, 0, 2}
	// v := binary.BigEndian.Uint32(i)
	// fmt.Println(v)

	// Assuming big-endian byte order, change to binary.LittleEndian if needed.
	i := []byte{0, 0, 0, 2}

	var intValue int32
	err := binary.Read(bytes.NewReader(i), binary.BigEndian, &intValue)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(intValue)

	// json.Unmarshal(b, v)

	fmt.Println(string(b))
}
