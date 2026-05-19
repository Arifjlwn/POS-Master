package models

import "time"

// FnBDetail menyimpan detail pesanan tiap menu (mirip struct laundry_detail)
type FnBDetail struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID uint      `json:"transaction_id"`
	ProductID     uint      `json:"product_id"`
	Qty           int       `json:"qty"`
	HargaSatuan   float64   `json:"harga_satuan"`
	SubTotal      float64   `json:"sub_total"`
	Notes         string    `gorm:"type:varchar(255)" json:"notes"` // 🚀 Catatan khusus koki (pedas, es dikit)
	StatusDapur   string    `gorm:"type:varchar(50);default:'PROSES'" json:"status_dapur"` // PROSES / SELESAI
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relasi ke tabel Core Global
	Transaction   Transaction `gorm:"foreignKey:TransactionID" json:"-"`
	Product       Product     `gorm:"foreignKey:ProductID" json:"product"`
}