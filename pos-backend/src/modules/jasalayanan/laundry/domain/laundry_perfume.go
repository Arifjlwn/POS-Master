package domain

import "time"

type Perfume struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoreID   uint      `gorm:"index;not null" json:"store_id"` 
	Nama      string    `gorm:"type:varchar(100);not null" json:"nama"` 
	Deskripsi string    `gorm:"type:varchar(255);default:''" json:"deskripsi"` 
	Harga     float64   `gorm:"type:decimal(10,2);default:0.00" json:"harga"` 
	Status    string    `gorm:"type:varchar(20);index;default:'AKTIF'" json:"status"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Perfume) TableName() string {
	return "laundry_perfumes"
}