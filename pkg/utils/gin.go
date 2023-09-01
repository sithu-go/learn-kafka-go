package utils

import "github.com/gin-gonic/gin"

func GetUserIDFromRequest(ctx *gin.Context) (string, error) {
	userID := ctx.Param("userID")
	if userID == "" {
		return "", ErrNoMessagesFound
	}
	return userID, nil
}
