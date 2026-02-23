package controllers

import (
	"account-management-system/backend/database"
	"account-management-system/backend/models"
	"account-management-system/backend/utils"

	"github.com/gin-gonic/gin"
)

// LogOperation 记录操作日志
func LogOperation(c *gin.Context, operationType, module, description string) {
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "unknown"
	}

	database.DB.Create(&models.OperationLog{
		LogID:         utils.GenerateID("log"),
		UserID:        userID,
		OperationType: operationType,
		Module:        module,
		Description:   description,
		RequestIP:     c.ClientIP(),
		Status:        1,
	})
}
