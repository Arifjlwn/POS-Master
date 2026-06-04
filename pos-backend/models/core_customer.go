package models

import "time"

type Customer struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex" json:"public_id"`

	StoreID uint `gorm:"not null;index" json:"store_id"`

	Nama string `gorm:"size:150;not null;index" json:"nama"`

	NoWhatsapp string `gorm:"size:20;index" json:"no_whatsapp"`
	Email       string `gorm:"size:150;index" json:"email"`

	Alamat string `gorm:"type:text" json:"alamat"`

	// Membership
	IsMember bool `gorm:"default:false" json:"is_member"`

	MemberCode string `gorm:"size:50;index" json:"member_code"`

	LoyaltyPoint int `gorm:"default:0" json:"loyalty_point"`

	// Statistik Customer
	TotalTransaction int     `gorm:"default:0" json:"total_transaction"`
	TotalSpending    float64 `gorm:"type:decimal(15,2);default:0" json:"total_spending"`

	LastTransactionAt *time.Time `json:"last_transaction_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Store Store `gorm:"foreignKey:StoreID" json:"-"`
}