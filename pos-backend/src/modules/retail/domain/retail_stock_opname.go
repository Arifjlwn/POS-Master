package domain

import (
	"time"
	"pos-backend/models"
)

// ==========================================
// 1. MODEL CONFIG UNIT: STOCK OPNAME REGULER
// ==========================================

type StockOpname struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PublicID  string    `gorm:"size:26;uniqueIndex;not null" json:"public_id"` // Masking ID eksternal via ULID

	StoreID   uint      `gorm:"index;not null" json:"store_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`

	Notes     string    `gorm:"type:text" json:"notes"`
	Status    string    `gorm:"type:varchar(50);default:'PENDING_APPROVAL';index" json:"status"` // PENDING_APPROVAL / APPROVED
	BuktiBar  string    `gorm:"type:text" json:"bukti_bar"` // URL Panjang Supabase aman tanpa takut truncate

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Details   []StockOpnameDetail `gorm:"foreignKey:OpnameID;constraint:OnDelete:CASCADE;" json:"details"` // 🚀 Ditambah constraint cascade biar aman  
}

func (StockOpname) TableName() string { return "retail_stock_opnames" }

type StockOpnameDetail struct {
	ID        uint           `gorm:"primaryKey" json:"id"`

	OpnameID  uint           `gorm:"index;not null" json:"opname_id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`

	SystemQty int            `gorm:"not null" json:"system_qty"`
	ActualQty int            `gorm:"not null" json:"actual_qty"`
	Selisih   int            `gorm:"not null" json:"selisih"`
	NilaiUang float64        `gorm:"type:decimal(18,2);not null;default:0" json:"nilai_uang"` // Nominal kerugian/keuntungan rupiah stok

	Alasan    string         `gorm:"type:text" json:"alasan"`
	Product   models.Product `gorm:"foreignKey:ProductID" json:"product"`
}

// 🛡️ FIX NAMA TABEL: Bersihkan typo double kata 'stock_stock'  biar query GORM lurus total!
func (StockOpnameDetail) TableName() string { return "retail_stock_opname_details" }

// ==========================================
// 2. MODEL CONFIG UNIT: KLAIM BARANG NYEMPIL (ADJUSTMENT)
// ==========================================

type StockAdjustment struct {
	ID        uint                    `gorm:"primaryKey" json:"id"`
	PublicID  string                  `gorm:"size:26;uniqueIndex;not null" json:"public_id"` // Masking ID eksternal via ULID

	StoreID   uint                    `gorm:"index;not null" json:"store_id"`
	UserID    uint                    `gorm:"not null;index" json:"user_id"`

	Notes     string                  `gorm:"type:text" json:"notes"` // Default: "Klaim Barang Nyempil"
	Status    string                  `gorm:"type:varchar(50);default:'PENDING_APPROVAL';index" json:"status"` // PENDING_APPROVAL / APPROVED
	BuktiBar  string                  `gorm:"type:text" json:"bukti_bar"` // URL Panjang Supabase aman

	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
	Details   []StockAdjustmentDetail `gorm:"foreignKey:AdjustmentID;constraint:OnDelete:CASCADE;" json:"details"`
}

func (StockAdjustment) TableName() string { return "retail_stock_adjustments" }

type StockAdjustmentDetail struct {
	ID           uint           `gorm:"primaryKey" json:"id"`

	AdjustmentID uint           `gorm:"index;not null" json:"adjustment_id"`
	ProductID    uint           `gorm:"not null;index" json:"product_id"`

	Qty          int            `gorm:"not null" json:"qty"` // Jumlah fisik tak sengaja ketemu
	Alasan       string         `gorm:"type:text" json:"alasan"`
	Product      models.Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (StockAdjustmentDetail) TableName() string { return "retail_stock_adjustment_details" }