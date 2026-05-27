package domain

import (
	"time"
	"pos-backend/models" // Import Product global
)

// ==========================================
// 1. MODEL CONFIG UNIT: STOCK OPNAME REGULER
// ==========================================

type StockOpname struct {
	ID        uint                `gorm:"primaryKey" json:"id"`
	StoreID   uint                `gorm:"index;not null" json:"store_id"`
	UserID    uint                `json:"user_id"` 
	Notes     string              `json:"notes"`   
	Status    string              `gorm:"type:varchar(50);default:'PENDING_APPROVAL'" json:"status"` // PENDING_APPROVAL / APPROVED
	CreatedAt time.Time           `json:"created_at"`
	Details   []StockOpnameDetail `gorm:"foreignKey:OpnameID" json:"details"`
}

func (StockOpname) TableName() string { return "retail_stock_opnames" }

type StockOpnameDetail struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OpnameID  uint           `gorm:"index;not null" json:"opname_id"`
	ProductID uint           `json:"product_id"`
	SystemQty int            `json:"system_qty"` 
	ActualQty int            `json:"actual_qty"` 
	Selisih   int            `json:"selisih"`    
	Alasan    string         `gorm:"type:text" json:"alasan"`
	Product   models.Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (StockOpnameDetail) TableName() string { return "retail_stock_opname_details" }

// ==========================================
// 2. MODEL CONFIG UNIT: KLAIM BARANG NYEMPIL
// ==========================================

type StockAdjustment struct {
	ID        uint                    `gorm:"primaryKey" json:"id"`
	StoreID   uint                    `gorm:"index;not null" json:"store_id"` // Biar aman per cabang toko
	UserID    uint                    `json:"user_id"`                        // Siapa kasir/karyawan yang nemu
	Notes     string                  `json:"notes"`                          // Default: "Klaim Barang Nyempil"
	Status    string                  `gorm:"type:varchar(50);default:'PENDING_APPROVAL'" json:"status"` // Owner wajib approve dulu baru stok nambah
	CreatedAt time.Time               `json:"current_time" json:"created_at"`
	Details   []StockAdjustmentDetail `gorm:"foreignKey:AdjustmentID" json:"details"`
}

func (StockAdjustment) TableName() string { return "retail_stock_adjustments" }

type StockAdjustmentDetail struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	AdjustmentID uint           `gorm:"index;not null" json:"adjustment_id"`
	ProductID    uint           `json:"product_id"`
	Qty          int            `json:"qty"`                     // Jumlah fisik yang ga sengaja ketemu
	Alasan       string         `gorm:"type:text" json:"alasan"` // Ketemu di mana / rak mana
	Product      models.Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (StockAdjustmentDetail) TableName() string { return "retail_stock_adjustment_details" }