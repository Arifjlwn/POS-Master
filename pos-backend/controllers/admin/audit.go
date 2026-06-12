package admin

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditController struct {
	DB *gorm.DB
}

// GetDetailedAuditLogs handles GET ALL AUDIT LOGS DENGAN LIVE SEARCH PUSAT
func (a *AuditController) GetDetailedAuditLogs(c *gin.Context) {
	search := c.Query("search")

	type DetailedLogResponse struct {
		ID        uint      `json:"id"`
		Action    string    `json:"action"`
		Details   string    `json:"details"`
		UserEmail string    `json:"user_email"`
		UserName  string    `json:"user_name"`
		IPAddress string    `json:"ip_address"`
		UserAgent string    `json:"user_agent"`
		CreatedAt time.Time `json:"created_at"`
	}

	var logs []DetailedLogResponse

	query := a.DB.Table("audit_logs").
		Select(`
			audit_logs.id, 
			audit_logs.action, 
			audit_logs.details, 
			audit_logs.user_email,
			COALESCE(users.name, audit_logs.user_email) as user_name, 
			audit_logs.ip_address,
			audit_logs.user_agent,
			audit_logs.created_at
		`).
		Joins("LEFT JOIN users ON users.id = audit_logs.user_id")

	if search != "" {
		searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(search))
		query = query.Where("LOWER(audit_logs.action) LIKE ? OR LOWER(audit_logs.details) LIKE ? OR LOWER(users.name) LIKE ? OR LOWER(audit_logs.user_email) LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm)
	}

	err := query.Order("audit_logs.created_at desc").Scan(&logs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik data master log audit pusat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   logs,
	})
}