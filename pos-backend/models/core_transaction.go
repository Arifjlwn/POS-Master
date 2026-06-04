package models

import "time"

// ======================================================
// TRANSACTION HEADER
// ======================================================

type Transaction struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex;not null" json:"public_id"`

	SessionID uint `gorm:"not null;index" json:"session_id"`

	StoreID uint `gorm:"not null;index" json:"store_id"`
	Store   Store `gorm:"foreignKey:StoreID" json:"store"`

	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	NoInvoice string `gorm:"type:varchar(50);not null;index" json:"no_invoice"`

	SubTotal   float64 `gorm:"type:decimal(18,2);not null" json:"sub_total"`
	Pajak      float64 `gorm:"type:decimal(18,2);default:0" json:"pajak"`
	Pembulatan float64 `gorm:"type:decimal(18,2);default:0" json:"pembulatan"`

	TotalHarga float64 `gorm:"type:decimal(18,2);not null" json:"total_harga"`

	MetodeBayar string `gorm:"type:varchar(20);not null;default:'Cash';index" json:"metode_bayar"`

	StatusBayar string `gorm:"type:varchar(20);not null;default:'LUNAS';index" json:"status_bayar"`

	TipeBisnis string `gorm:"type:varchar(20);not null;default:'RETAIL';index" json:"tipe_bisnis"`

	StatusPesanan string `gorm:"type:varchar(50);default:'SELESAI';index" json:"status_pesanan"`

	BuktiTransfer string `gorm:"type:text" json:"bukti_transfer"`

	NominalBayar float64 `gorm:"type:decimal(18,2);not null" json:"nominal_bayar"`

	Kembalian float64 `gorm:"type:decimal(18,2);not null" json:"kembalian"`

	Details []TransactionDetail `gorm:"foreignKey:TransactionID" json:"details"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ======================================================
// TRANSACTION DETAIL
// ======================================================

type TransactionDetail struct {
	ID uint `gorm:"primaryKey" json:"id"`

	TransactionID uint `gorm:"not null;index" json:"transaction_id"`

	ProductID uint `gorm:"not null;index" json:"product_id"`

	Product Product `gorm:"foreignKey:ProductID" json:"product"`

	// Snapshot data
	NamaProdukSnapshot string `gorm:"type:varchar(150);not null" json:"nama_produk_snapshot"`

	SKUProductSnapshot string `gorm:"type:varchar(50)" json:"sku_product_snapshot"`

	HargaSatuan float64 `gorm:"type:decimal(18,2);not null" json:"harga_satuan"`

	Kuantitas int `gorm:"not null" json:"kuantitas"`

	ItemType string `gorm:"type:varchar(30);default:'PRODUCT';index" json:"item_type"`

	DetailNotes string `gorm:"type:text" json:"detail_notes"`

	SubTotal float64 `gorm:"type:decimal(18,2);not null" json:"sub_total"`
}