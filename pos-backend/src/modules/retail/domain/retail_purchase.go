package domain

import (
	"pos-backend/models"
	"time"
)

// ==========================================
// 🚀 MODEL CONFIG UNIT: PURCHASE / LPB FAKTUR
// ==========================================

type Purchase struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex;not null" json:"public_id"` // Masking ID eksternal via ULID

	StoreID uint `gorm:"index;not null" json:"store_id"`
	UserID  uint `gorm:"not null" json:"user_id"`

	SupplierName string  `gorm:"type:varchar(100)" json:"supplier_name"`
	NoFaktur     string  `gorm:"type:varchar(50);index" json:"no_faktur"`
	TotalItem    int     `json:"total_item"`
	TotalHarga   float64 `gorm:"type:decimal(18,2);not null;default:0" json:"total_harga"`            // Untuk audit mutasi uang keluar
	StatusBayar  string  `gorm:"type:varchar(20);not null;default:'LUNAS';index" json:"status_bayar"` // LUNAS / HUTANG

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 🚀 FIX SECURITY CONSTRAINT: Tambah OnDelete:CASCADE biar DB ga nyimpen data sampah/yatim piatu !
	Details []PurchaseDetail `gorm:"foreignKey:PurchaseID;constraint:OnDelete:CASCADE;" json:"details"`
}

func (Purchase) TableName() string { return "retail_purchases" }

type PurchaseDetail struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// HYBRID OPTIMAL: PublicID ULID dibuang dari detail, hemat index space B-Tree !

	PurchaseID uint `gorm:"index;not null" json:"purchase_id"`
	ProductID  uint `gorm:"not null;index" json:"product_id"`

	QtyMasuk int `gorm:"not null" json:"qty_masuk"`

	// 🛡️ RE-CALIBRATION JSON TAG: Tetap gorm harga_beli tapi JSON disamakan dengan standard payload
	HargaModal float64 `gorm:"type:decimal(18,2);not null;default:0" json:"harga_modal"`
	SubTotal   float64 `gorm:"type:decimal(18,2);not null;default:0" json:"sub_total"` // Auto-calculation per row item server-side

	Product models.Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (PurchaseDetail) TableName() string { return "retail_purchase_details" }

// ==========================================
// 📈 STRUCT DATA COMPLIANCE REPORT ANALYTICS
// ==========================================

type BestSeller struct {
	NamaProduk  string  `json:"nama_produk"`
	SKU         string  `json:"sku"`
	QtyTerjual  int     `json:"qty_terjual"`
	TotalOmzet  float64 `json:"total_omzet"`
	SatuanDasar string  `json:"satuan_dasar"`
}

type InboundItemResponse struct {
	ID         uint    `json:"id"`
	NamaProduk string  `json:"nama_produk"`
	Qty        int     `json:"qty"`
	HargaModal float64 `json:"harga_modal"`
	SubTotal   float64 `json:"sub_total"`
}

type InboundReportResponse struct {
	ID           uint                  `json:"id"`
	NoFaktur     string                `json:"no_faktur"`
	NamaSupplier string                `json:"nama_supplier"`
	TotalItem    int                   `json:"total_item"`
	TotalModal   float64               `json:"total_modal"`
	CreatedAt    time.Time             `json:"created_at"`
	Items        []InboundItemResponse `json:"items"`
}
