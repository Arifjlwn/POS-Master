package models

import "time"

// 1. Model untuk Kepala Struk (Header)
type Transaction struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    SessionID    uint      `gorm:"not null" json:"session_id"` // 🚀 KUNCI UTAMA UNTUK CLOSING
    StoreID      uint      `gorm:"not null" json:"store_id"`
    UserID       uint      `gorm:"not null" json:"user_id"`
    NoInvoice    string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"no_invoice"`

    SubTotal     float64   `gorm:"type:decimal(10,2);not null" json:"sub_total"`
    Pajak        float64   `gorm:"type:decimal(10,2);default:0" json:"pajak"`
    Pembulatan   float64   `gorm:"type:decimal(10,2);default:0" json:"pembulatan"`
    TotalHarga   float64   `gorm:"type:decimal(10,2);not null" json:"total_harga"`
    
    MetodeBayar  string    `gorm:"type:varchar(20);not null;default:'Cash'" json:"metode_bayar"` // 🚀 TUNAI / QRIS
    
    NominalBayar float64   `gorm:"type:decimal(10,2);not null" json:"nominal_bayar"`
    Kembalian    float64   `gorm:"type:decimal(10,2);not null" json:"kembalian"`
    CreatedAt    time.Time `json:"created_at"`

    Details []TransactionDetail `gorm:"foreignKey:TransactionID" json:"details"`
    User    User              `gorm:"foreignKey:UserID" json:"User"`
}

// 2. Model untuk Rincian Barang di Struk (Body)
type TransactionDetail struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	TransactionID uint    `gorm:"not null" json:"transaction_id"`
	ProductID     uint    `gorm:"not null" json:"product_id"`
	HargaSatuan   float64 `gorm:"type:decimal(10,2);not null" json:"harga_satuan"`
	Kuantitas     int     `gorm:"not null" json:"kuantitas"`
	SubTotal      float64 `gorm:"type:decimal(10,2);not null" json:"sub_total"`

	// Relasi untuk narik nama produk (opsional, buat nampilin di nota)
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}