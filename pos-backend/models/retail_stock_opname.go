package models

import (
	"time"
)

type StockOpname struct {
	ID        uint                `gorm:"primaryKey" json:"id"`
	StoreID   uint                `json:"store_id"`
	UserID    uint                `json:"user_id"` // Siapa yang hitung
	Notes     string              `json:"notes"`   // Misal: "SO Bulanan"
	CreatedAt time.Time           `json:"created_at"`
	Details   []StockOpnameDetail `gorm:"foreignKey:OpnameID" json:"details"`
}

type StockOpnameDetail struct {
	ID        uint    `gorm:"primaryKey"`
	OpnameID  uint    `json:"opname_id"`
	ProductID uint    `json:"product_id"`
	SystemQty int     `json:"system_qty"` // Angka sistem saat SO dikunci
	ActualQty int     `json:"actual_qty"` // Angka yang diinput karyawan
	Selisih   int     `json:"selisih"`    // (Actual - System)
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
}