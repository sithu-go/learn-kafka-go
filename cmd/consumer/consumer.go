package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"learn-kafka/pkg/models"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

const (
	ConsumerGroup = "notifications-group"
	ConsumerTopic = "notifications"
	ConsumerPort  = ":8889"
)

var (
	KafkaServerAddresses = []string{"localhost:9991", "localhost:9992", "localhost:9993"}
)

// ============== HELPER FUNCTIONS ==============
var ErrNoMessagesFound = errors.New("no messages found")

func getUserIDFromRequest(ctx *gin.Context) (string, error) {
	userID := ctx.Param("userID")
	if userID == "" {
		return "", ErrNoMessagesFound
	}
	return userID, nil
}

// ====== NOTIFICATION STORAGE ======
type UserNotifications map[int32][]models.Notification

type NotificationStore struct {
	data UserNotifications
	mu   sync.RWMutex
}

func (ns *NotificationStore) Add(userID int32, notification models.Notification) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.data[userID] = append(ns.data[userID], notification)
}

func (ns *NotificationStore) Get(userID int32) []models.Notification {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	return ns.data[userID]
}

// ============== KAFKA RELATED FUNCTIONS ==============
type Consumer struct {
	store *NotificationStore
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		var userID int32
		// converting binary into int32
		if err := binary.Read(bytes.NewReader(msg.Key), binary.BigEndian, &userID); err != nil {
			fmt.Println("Err", err)
			k, err := strconv.Atoi(string(msg.Key))
			if err != nil {
				fmt.Println("Err converting string into int", err)
			} else {
				userID = int32(k)
			}
		}

		var notification models.Notification
		err := json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}
		fmt.Println("Message Received", msg.Key, userID, "SS")
		// fmt.Printf("%+v\n", msg)
		// fmt.Printf("%+v\n", notification)

		consumer.store.Add(userID, notification)

		fmt.Printf("%+v\n", consumer.store.data)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func initializeConsumerGroup() (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.ClientID = "strange-g"

	consumerGroup, err := sarama.NewConsumerGroup(KafkaServerAddresses, ConsumerGroup, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize consumer group: %w", err)
	}

	return consumerGroup, nil
}

func setupConsumerGroup(ctx context.Context, store *NotificationStore) {
	consumerGroup, err := initializeConsumerGroup()
	if err != nil {
		log.Printf("initialization error: %v", err)
	}
	defer consumerGroup.Close()

	consumer := &Consumer{
		store: store,
	}

	for {
		err = consumerGroup.Consume(ctx, []string{ConsumerTopic}, consumer)
		if err != nil {
			log.Printf("error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func handleNotifications(ctx *gin.Context, store *NotificationStore) {
	userID, err := getUserIDFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Println(err)
	}

	notes := store.Get(int32(uid))
	if len(notes) == 0 {
		ctx.JSON(http.StatusOK,
			gin.H{
				"message":       "No notifications found for user",
				"notifications": []models.Notification{},
			})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications": notes})
}

func main() {
	store := &NotificationStore{
		data: make(UserNotifications),
	}

	ctx, cancel := context.WithCancel(context.Background())
	// beaware that
	go setupConsumerGroup(ctx, store)
	defer cancel()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/notifications/:userID", func(ctx *gin.Context) {
		handleNotifications(ctx, store)
	})

	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 "+
		"started at http://localhost%s\n", ConsumerGroup, ConsumerPort)

	if err := router.Run(ConsumerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
