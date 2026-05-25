package models

import "time"

// 1. Model untuk Kepala Struk Global (Header)
type Transaction struct {
	ID            uint                `gorm:"primaryKey" json:"id"`
	SessionID     uint                `gorm:"not null" json:"session_id"` // 🚀 KUNCI UTAMA UNTUK CLOSING
	StoreID       uint                `gorm:"not null" json:"store_id"`
	Store         Store               `gorm:"foreignKey:StoreID" json:"Store"`
	UserID        uint                `gorm:"not null" json:"user_id"`
	NoInvoice     string              `gorm:"type:varchar(50);uniqueIndex;not null" json:"no_invoice"`

	SubTotal      float64             `gorm:"type:decimal(10,2);not null" json:"sub_total"`
	Pajak         float64             `gorm:"type:decimal(10,2);default:0" json:"pajak"`
	Pembulatan    float64             `gorm:"type:decimal(10,2);default:0" json:"pembulatan"`
	TotalHarga    float64             `gorm:"type:decimal(10,2);not null" json:"total_harga"`
	
	MetodeBayar   string              `gorm:"type:varchar(20);not null;default:'Cash'" json:"metode_bayar"`
	StatusBayar   string              `gorm:"type:varchar(20);not null;default:'LUNAS'" json:"status_bayar"`
	
	// 🚀 GLOBAL FITUR: Default disesuaikan lewat input controller/frontend
	TipeBisnis    string              `gorm:"type:varchar(20);not null;default:'RETAIL'" json:"tipe_bisnis"`   // RETAIL, LAUNDRY, FNB
	StatusPesanan string              `gorm:"type:varchar(50);default:'SELESAI'" json:"status_pesanan"`       // SELESAI, ANTRI, PROSES
	
	BuktiTransfer string              `gorm:"type:varchar(255)" json:"bukti_transfer"`
	NominalBayar  float64             `gorm:"type:decimal(10,2);not null" json:"nominal_bayar"`
	Kembalian     float64             `gorm:"type:decimal(10,2);not null" json:"kembalian"`
	CreatedAt     time.Time           `json:"created_at"`

	Details       []TransactionDetail `gorm:"foreignKey:TransactionID" json:"details"`
	User          User                `gorm:"foreignKey:UserID" json:"User"`
}

// 2. Model untuk Rincian Barang di Struk Global (Body)
type TransactionDetail struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	TransactionID uint    `gorm:"not null" json:"transaction_id"`
	ProductID     uint    `gorm:"not null" json:"product_id"`
	HargaSatuan   float64 `gorm:"type:decimal(10,2);not null" json:"harga_satuan"`
	Kuantitas     int     `gorm:"not null" json:"kuantitas"`

	// 🚀 TRICK ARCHITECTURE: Hapus NamaParfum & HargaParfum, ganti dengan 2 kolom sakti ini:
	ItemType      string  `gorm:"type:varchar(30);default:'PRODUCT'" json:"item_type"` // PRODUCT, SERVICE, MENU
	DetailNotes   string  `gorm:"type:text" json:"detail_notes"`                       // Buat nyimpen catatan (ex: "Parfum: Premium, Wangi Lily" atau "F&B: Level 5")
	
	SubTotal      float64 `gorm:"type:decimal(10,2);not null" json:"sub_total"`

	// Relasi untuk narik nama produk
	Product       Product `gorm:"foreignKey:ProductID" json:"product"`
}