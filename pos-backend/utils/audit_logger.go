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

func RecordSystemLog(c *gin.Context, action string, targetULID string, details string) error {
	logData := models.AuditLog{
		UserID:         0,             // 0 menandakan ini dieksekusi oleh Sistem/Publik
		UserEmail:      "SYSTEM_AUTO", // Nama user yang bakal nongol di Mission Control
		Action:         action,
		TargetPublicID: targetULID,
		Details:        details,
		IPAddress:      c.ClientIP(),
		UserAgent:      c.Request.UserAgent(),
	}

	return src.DB.Create(&logData).Error
}
