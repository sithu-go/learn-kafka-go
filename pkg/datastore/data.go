package datastore

import "learn-kafka/pkg/models"

// ====== NOTIFICATION STORAGE ======
type UserNotifications map[string][]models.Notification

var (
	Users = []models.User{
		{ID: 1, Name: "Emma"},
		{ID: 2, Name: "Bruno"},
		{ID: 3, Name: "Rick"},
		{ID: 4, Name: "Lena"},
	}
)
