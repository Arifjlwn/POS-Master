package models

import "time"

// Header Penerimaan Barang
type Purchase struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	StoreID      uint             `gorm:"not null" json:"store_id"`
	UserID       uint             `gorm:"not null" json:"user_id"` // Siapa yang input
	SupplierName string           `gorm:"type:varchar(100)" json:"supplier_name"`
	NoFaktur     string           `gorm:"type:varchar(50)" json:"no_faktur"`
	TotalItem    int              `json:"total_item"`
	CreatedAt    time.Time        `json:"created_at"`
	Details      []PurchaseDetail `gorm:"foreignKey:PurchaseID" json:"details"`
}

// Rincian Barang yang Masuk
type PurchaseDetail struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	PurchaseID uint    `gorm:"not null" json:"purchase_id"`
	ProductID  uint    `gorm:"not null" json:"product_id"`
	QtyMasuk   int     `gorm:"not null" json:"qty_masuk"`
	HargaModal  float64 `gorm:"type:decimal(10,2)" json:"harga_beli"` // Harga dari supplier
	Product    Product `gorm:"foreignKey:ProductID" json:"product"`
}