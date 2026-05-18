package models

import (
	"time"
)

// Store merepresentasikan tabel 'stores' di Supabase yang sinkron dengan Form Setup
type Store struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	NamaToko     string    `gorm:"type:varchar(100);not null" json:"nama_toko"`
	Telepon      string    `gorm:"type:varchar(20)" json:"telepon"` // No WhatsApp Bisnis (+62...)
	
	// 🚀 EKOSISTEM UTAMA SAAS
	// Berisi string gabungan kategori dan spesifikasi bisnis. Contoh: "Jasa - Laundry"
	BusinessType string    `gorm:"type:varchar(100);column:business_type;not null" json:"tipe_bisnis"` 
	
	// Menyimpan daftar modul pilihan premium dari nomor 3 (contoh: "kasir,whatsapp,absensi")
	FiturAktif   string    `gorm:"type:text;column:fitur_aktif" json:"fitur_aktif"` 

	// 🚀 PECAHAN LOKASI OPERASIONAL (SINKRON NOMOR 2)
	Alamat       string    `gorm:"type:text" json:"alamat"` // Gabungan full alamat jalan text untuk kemudahan cetak nota
	Provinsi     string    `gorm:"type:varchar(100)" json:"provinsi"`
	Kota         string    `gorm:"type:varchar(100)" json:"kota"`
	Kecamatan    string    `gorm:"type:varchar(100)" json:"kecamatan"`
	Kelurahan    string    `gorm:"type:varchar(100)" json:"kelurahan"`
	KodePos      string    `gorm:"type:varchar(10)" json:"kode_pos"`

	PajakPersen  float64   `gorm:"type:decimal(5,2);default:0" json:"pajak_persen"` 
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	PaymentType   string `gorm:"type:varchar(20);default:'PRIBADI'" json:"payment_type"` // 'PRIBADI' atau 'GATEWAY'
	QrisImage     string `gorm:"type:varchar(255)" json:"qris_image"`                     // Path file barcode QRIS
	ReceiptFooter string `gorm:"type:text" json:"receipt_footer"`                         // Catatan kaki di struk thermal
}