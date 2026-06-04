package domain

import (
	"time"
	"pos-backend/models"
)

type ProductReturn struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PublicID  string         `gorm:"size:26;uniqueIndex;not null" json:"public_id"` // 🚀 Wajib ULID karena tabelnya model FLAT

	ReturnNo  string         `gorm:"type:varchar(50);index;not null" json:"return_no"`
	StoreID   uint           `gorm:"index;not null" json:"store_id"`

	ProductID uint           `gorm:"not null;index" json:"product_id"`
	Product   models.Product `gorm:"foreignKey:ProductID" json:"product"`

	UserID    uint           `gorm:"not null;index" json:"user_id"`
	User      models.User    `gorm:"foreignKey:UserID" json:"user"`

	Qty       int            `gorm:"not null;default:0" json:"qty"`
	Alasan    string         `gorm:"type:varchar(255)" json:"alasan"`
	Catatan   string         `gorm:"type:text" json:"catatan"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (ProductReturn) TableName() string { return "retail_product_returns" }