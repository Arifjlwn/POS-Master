package utils

import (
	"pos-backend/models"
	src "pos-backend/src/core/config"
	"github.com/gin-gonic/gin"
)

// RecordAdminAction otomatis menangkap IP & browser admin untuk dicatat ke database
func RecordAdminAction(c *gin.Context, action string, targetULID string, details string) error {
	userIDRaw, _ := c.Get("user_id")
	userID := uint(userIDRaw.(float64))
	
	emailRaw, _ := c.Get("email")
	userEmail := "unknown_admin"
	if emailRaw != nil {
		userEmail = emailRaw.(string)
	}

	logData := models.AuditLog{
		UserID:         userID,
		UserEmail:      userEmail,
		Action:         action,
		TargetPublicID: targetULID,
		Details:        details,
		IPAddress:      c.ClientIP(),
		UserAgent:      c.Request.UserAgent(),
	}

	return src.DB.Create(&logData).Error
}