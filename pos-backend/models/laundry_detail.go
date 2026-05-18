package models

import "time"

// TransactionLaundryDetail menyimpan detail khusus untuk transaksi jasa/laundry
type TransactionLaundryDetail struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID uint      `json:"transaction_id"` // Nembak ke tabel transactions utama
	ProductID     uint      `json:"product_id"`     // Nembak ke tabel products (Layanan Cuci)
	NamaPelanggan string    `gorm:"type:varchar(100)" json:"nama_pelanggan"`
	NoWhatsapp    string    `gorm:"type:varchar(20)" json:"no_whatsapp"`
	BeratKg       float64   `json:"berat_kg"`       // Bisa desimal misal 2.5 Kg
	HargaPerKg    float64   `json:"harga_per_kg"`   // Harga jasa saat transaksi
	SubTotal      float64   `json:"sub_total"`
	FotoBarang string `gorm:"type:varchar(255)" json:"foto_barang"`
	StatusCucian  string    `gorm:"type:varchar(50);default:'ANTRI'" json:"status_cucian"` // ANTRI, DICUCI, SELESAI, DIAMBIL
	EstimasiWaktu time.Time `json:"estimasi_waktu"` // Kapan kira-kira kelar
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}