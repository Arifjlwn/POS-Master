package models

import "time"

type Perfume struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoreID   uint      `gorm:"not null" json:"store_id"`
	Nama      string    `gorm:"type:varchar(50);not null" json:"nama"`
	Harga     float64   `gorm:"type:decimal(10,2);default:0" json:"harga"` // Charge tambahan
	Status    string    `gorm:"type:varchar(20);default:'Tersedia'" json:"status"` // Tersedia / Habis
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Store     Store     `gorm:"foreignKey:StoreID" json:"-"`
}