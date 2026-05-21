package domain

import "time"

type Perfume struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoreID   uint      `gorm:"index;not null" json:"store_id"` // Index biar filter parfum per toko cepet
	Nama      string    `gorm:"type:varchar(50);not null" json:"nama"`
	Harga     float64   `gorm:"type:decimal(10,2);default:0" json:"harga"` 
	Status    string    `gorm:"type:varchar(20);default:'Tersedia'" json:"status"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Perfume) TableName() string {
	return "laundry_perfumes"
}