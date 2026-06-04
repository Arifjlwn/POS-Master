package models

import "time"

type ShiftClosing struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex" json:"public_id"`

	StoreID   uint `gorm:"not null;index" json:"store_id"`
	SessionID uint `gorm:"not null;index" json:"session_id"`
	UserID    uint `gorm:"not null;index" json:"user_id"`

	// Waktu Shift
	OpenedAt time.Time `json:"opened_at"`
	ClosedAt time.Time `json:"closed_at"`

	// Ringkasan Penjualan
	TotalTransaction int `gorm:"default:0" json:"total_transaction"`

	NetSales      float64 `gorm:"type:decimal(15,2);default:0" json:"net_sales"`
	TotalTax      float64 `gorm:"type:decimal(15,2);default:0" json:"total_tax"`
	SalesCash     float64 `gorm:"type:decimal(15,2);default:0" json:"sales_cash"`
	SalesNonCash  float64 `gorm:"type:decimal(15,2);default:0" json:"sales_non_cash"`

	// Rekap Kas
	ModalAwal     float64 `gorm:"type:decimal(15,2);default:0" json:"modal_awal"`
	TotalExpected float64 `gorm:"type:decimal(15,2);default:0" json:"total_expected"`
	TotalActual   float64 `gorm:"type:decimal(15,2);default:0" json:"total_actual"`
	Selisih       float64 `gorm:"type:decimal(15,2);default:0" json:"selisih"`

	// Catatan Closing
	ClosingNote string `gorm:"type:text" json:"closing_note"`

	CreatedAt time.Time `json:"created_at"`

	// Relations
	User    User           `gorm:"foreignKey:UserID" json:"user"`
	Session CashierSession `gorm:"foreignKey:SessionID" json:"session"`
	Store   Store          `gorm:"foreignKey:StoreID" json:"store"`
}