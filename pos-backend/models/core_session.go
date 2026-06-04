package models

import "time"

type CashierSession struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex" json:"public_id"`

	StoreID uint `gorm:"not null;index" json:"store_id"`
	UserID  uint `gorm:"not null;index" json:"user_id"`

	// Kasir ke berapa
	StationNumber string `gorm:"size:20;index" json:"station_number"`

	// Modal Awal Shift
	ModalAwal float64 `gorm:"type:decimal(15,2);not null;default:0" json:"modal_awal"`

	// Rekap Kas
	TotalMasuk  float64 `gorm:"type:decimal(15,2);default:0" json:"total_masuk"`
	TotalKeluar float64 `gorm:"type:decimal(15,2);default:0" json:"total_keluar"`

	// Closing
	TotalAktual float64 `gorm:"type:decimal(15,2);default:0" json:"total_aktual"`
	Selisih     float64 `gorm:"type:decimal(15,2);default:0" json:"selisih"`

	// Catatan jika ada selisih
	ClosingNote string `gorm:"type:text" json:"closing_note"`

	// Status
	Status string `gorm:"size:20;index;default:'OPEN'" json:"status"`
	// OPEN
	// CLOSED
	// CANCELLED

	OpenedAt time.Time  `gorm:"index" json:"opened_at"`
	ClosedAt *time.Time `gorm:"index" json:"closed_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User  User  `gorm:"foreignKey:UserID" json:"user"`
	Store Store `gorm:"foreignKey:StoreID" json:"store"`

	Transactions []Transaction `gorm:"foreignKey:SessionID" json:"-"`
}