package domain

import (
	"time"
	"pos-backend/models" // Import Product global
)

type StockOpname struct {
	ID        uint                `gorm:"primaryKey" json:"id"`
	StoreID   uint                `gorm:"index;not null" json:"store_id"`
	UserID    uint                `json:"user_id"` 
	Notes     string              `json:"notes"`   
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
	Product   models.Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (StockOpnameDetail) TableName() string { return "retail_stock_opname_details" }