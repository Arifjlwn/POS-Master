package domain

import (
	"time"
	"pos-backend/models" // Import Product global
)

type Purchase struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	StoreID      uint             `gorm:"index;not null" json:"store_id"`
	UserID       uint             `gorm:"not null" json:"user_id"` 
	SupplierName string           `gorm:"type:varchar(100)" json:"supplier_name"`
	NoFaktur     string           `gorm:"type:varchar(50);index" json:"no_faktur"`
	TotalItem    int              `json:"total_item"`
	CreatedAt    time.Time        `json:"created_at"`
	Details      []PurchaseDetail `gorm:"foreignKey:PurchaseID" json:"details"`
}

func (Purchase) TableName() string { return "retail_purchases" }

type PurchaseDetail struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PurchaseID uint           `gorm:"index;not null" json:"purchase_id"`
	ProductID  uint           `gorm:"not null" json:"product_id"`
	QtyMasuk   int            `gorm:"not null" json:"qty_masuk"`
	HargaModal float64        `gorm:"type:decimal(10,2)" json:"harga_beli"` 
	Product    models.Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (PurchaseDetail) TableName() string { return "retail_purchase_details" }

type BestSeller struct {
	NamaProduk string  `json:"nama_produk"`
	SKU        string  `json:"sku"`
	QtyTerjual int     `json:"qty_terjual"`
	TotalOmzet float64 `json:"total_omzet"`
}