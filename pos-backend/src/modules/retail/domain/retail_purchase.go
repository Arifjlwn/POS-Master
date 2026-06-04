package domain

import (
	"time"
	"pos-backend/models"
)

type Purchase struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	PublicID     string           `gorm:"size:26;uniqueIndex;not null" json:"public_id"` // 🚀 AMAN: Masking ID eksternal via ULID

	StoreID      uint             `gorm:"index;not null" json:"store_id"`
	UserID       uint             `gorm:"not null" json:"user_id"` 

	SupplierName string           `gorm:"type:varchar(100)" json:"supplier_name"`
	NoFaktur     string           `gorm:"type:varchar(50);index" json:"no_faktur"`
	TotalItem    int              `json:"total_item"`
	TotalHarga   float64          `gorm:"type:decimal(18,2);not null;default:0" json:"total_harga"` // Untuk audit mutasi uang keluar
	StatusBayar  string           `gorm:"type:varchar(20);not null;default:'LUNAS';index" json:"status_bayar"` // LUNAS / HUTANG

	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Details      []PurchaseDetail `gorm:"foreignKey:PurchaseID" json:"details"`
}

func (Purchase) TableName() string { return "retail_purchases" }

type PurchaseDetail struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	// 🚀 HYBRID OPTIMAL: PublicID ULID dibuang biar database gak megap-megap nyimpen indeks string di tabel detail!

	PurchaseID uint           `gorm:"index;not null" json:"purchase_id"`
	ProductID  uint           `gorm:"not null" json:"product_id"`

	QtyMasuk   int            `gorm:"not null" json:"qty_masuk"`
	HargaModal float64        `gorm:"type:decimal(18,2);not null;default:0" json:"harga_beli"` 
	SubTotal   float64        `gorm:"type:decimal(18,2);not null;default:0" json:"sub_total"` // Auto-calculation per row item

	Product    models.Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (PurchaseDetail) TableName() string { return "retail_purchase_details" }

type BestSeller struct {
	NamaProduk  string  `json:"nama_produk"`
	SKU         string  `json:"sku"`
	QtyTerjual  int     `json:"qty_terjual"`
	TotalOmzet  float64 `json:"total_omzet"`
	SatuanDasar string  `json:"satuan_dasar"`
}