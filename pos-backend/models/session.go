package models

import (
	"time"
)

type CashierSession struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	StoreID       uint       `gorm:"not null" json:"store_id"`
	UserID        uint       `gorm:"not null" json:"user_id"`
	StationNumber string     `gorm:"size:10;not null" json:"station_number"` // Contoh: 01, 02
	ModalAwal     float64    `gorm:"type:decimal(15,2);not null" json:"modal_awal"`
	TotalMasuk    float64    `gorm:"type:decimal(15,2);default:0" json:"total_masuk"`
	TotalKeluar   float64    `gorm:"type:decimal(15,2);default:0" json:"total_keluar"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       *time.Time `json:"end_time"` // Pakai pointer biar bisa NULL (berarti belum closing)
	Status        string     `gorm:"size:20;default:'open'" json:"status"` // 'open' atau 'closed'
	
	// Relasi (Opsional tapi bagus buat laporan)
	User  User  `gorm:"foreignKey:UserID" json:"user"`
	Store Store `gorm:"foreignKey:StoreID" json:"store"`
}