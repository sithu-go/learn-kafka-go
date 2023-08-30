package main

import (
	"errors"
	"fmt"
	"learn-kafka/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	ProducerPort       = ":8080"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notifications"
)

var ErrUserNotFoundInProducer = errors.New("user not found")

// ============== HELPER FUNCTIONS ==============
func findUserByID(id int, users []models.User) (models.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFoundInProducer
}

func getIDFromRequest(formValue string, ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.PostForm(formValue))
	if err != nil {
		return 0, fmt.Errorf("failed to parse ID from form value %s: %w", formValue, err)
	}

	return id, nil
}

// // ============== KAFKA RELATED FUNCTIONS ==============
// func sendKafkaMessage(producer sarama.SyncProducer, users []models.User, ctx *gin.Context, fromID, toID int) error {

// 	message := ctx.PostForm("message")

// 	fromUser, err := findUserByID(fromID, users)
// 	if err != nil {
// 		return err
// 	}

// 	toUser, err := findUserByID(toID, users)
// 	if err != nil {
// 		return err
// 	}

// 	notification := models.Notification{
// 		From:    fromUser,
// 		To:      toUser,
// 		Message: message,
// 	}

// 	notificationJSON, err := json.Marshal(notification)
// 	if err != nil {
// 		return fmt.Errorf("failed to marsahl notification: %w", err)
// 	}

// 	msg := &sarama.ProducerMessage{
// 		Topic:     KafkaTopic,
// 		Key:       nil,
// 		Value:     nil,
// 		Headers:   []sarama.RecordHeader{},
// 		Metadata:  nil,
// 		Offset:    0,
// 		Partition: 0,
// 		Timestamp: time.Time{},
// 	}

// 	return nil
// }
