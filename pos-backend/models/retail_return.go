package models

import "time"

type ProductReturn struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ReturnNo  string    `json:"return_no"`
	StoreID   uint      `json:"store_id"`
	ProductID uint      `json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
	UserID    uint      `json:"user_id"`                             
	User      User      `gorm:"foreignKey:UserID" json:"user"`       
	Qty       int       `json:"qty"`
	Alasan    string    `json:"alasan"`  
	Catatan   string    `json:"catatan"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}