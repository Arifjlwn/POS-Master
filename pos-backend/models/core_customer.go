package models

import "time"

type Customer struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	StoreID    uint      `gorm:"not null" json:"store_id"`
	Nama       string    `gorm:"type:varchar(100);not null" json:"nama"`
	NoWhatsapp string    `gorm:"type:varchar(20);not null" json:"no_whatsapp"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}