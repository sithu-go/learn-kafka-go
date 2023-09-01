package utils

import (
	"learn-kafka/pkg/datastore"
	"learn-kafka/pkg/models"
)

// ============== HELPER FUNCTIONS ==============
func FindUserByID(id int) (models.User, error) {
	for _, user := range datastore.Users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFoundInProducer
}
