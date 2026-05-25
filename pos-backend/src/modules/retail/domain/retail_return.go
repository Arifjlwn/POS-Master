package domain

import (
	"time"
	"pos-backend/models" // Import Product & User global
)

type ProductReturn struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ReturnNo  string         `gorm:"type:varchar(50);index" json:"return_no"`
	StoreID   uint           `gorm:"index;not null" json:"store_id"`
	ProductID uint           `json:"product_id"`
	Product   models.Product `gorm:"foreignKey:ProductID" json:"product"`
	UserID    uint           `json:"user_id"`
	User      models.User    `gorm:"foreignKey:UserID" json:"user"`
	Qty       int            `json:"qty"`
	Alasan    string         `json:"alasan"`
	Catatan   string         `json:"catatan"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (ProductReturn) TableName() string { return "retail_product_returns" }