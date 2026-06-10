package models

import (
	"time"
)

// AuditLog mencatat tindakan sensitif yang dilakukan oleh Root Admin
type AuditLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`                  // Internal Primary Key
	UserID         uint      `gorm:"index" json:"user_id"`                  // ID Admin yang mengeksekusi
	UserEmail      string    `gorm:"type:varchar(150)" json:"user_email"`
	Action         string    `gorm:"type:varchar(100);index" json:"action"` // Contoh: "TENANT_SUSPEND"
	TargetPublicID string    `gorm:"type:varchar(100);index" json:"target_public_id"` // Target berbasis ULID
	Details        string    `gorm:"type:text" json:"details"`              // JSON info detail aksi
	IPAddress      string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent      string    `gorm:"type:text" json:"user_agent"`
	CreatedAt      time.Time `json:"created_at"`
}