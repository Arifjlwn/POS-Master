package models

import (
	"time"
)

// Store merepresentasikan tabel 'stores' di Supabase yang sinkron dengan Form Setup
type Store struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	NamaToko     string    `gorm:"type:varchar(100);not null" json:"nama_toko"`
	Telepon      string    `gorm:"type:varchar(20)" json:"telepon"`
	BusinessType string    `gorm:"type:varchar(100);column:business_type;not null" json:"business_type"` 
	FiturAktif   string    `gorm:"type:text;column:fitur_aktif" json:"fitur_aktif"` 

	Alamat       string    `gorm:"type:text" json:"alamat"`
	Provinsi     string    `gorm:"type:varchar(100)" json:"provinsi"`
	Kota         string    `gorm:"type:varchar(100)" json:"kota"`
	Kecamatan    string    `gorm:"type:varchar(100)" json:"kecamatan"`
	Kelurahan    string    `gorm:"type:varchar(100)" json:"kelurahan"`
	KodePos      string    `gorm:"type:varchar(10)" json:"kode_pos"`

	LogoURL      string    `gorm:"type:text" json:"logo_url"`

	IsTaxActive  bool      `gorm:"default:false" json:"is_tax_active"` 
    PajakPersen  float64   `gorm:"type:decimal(5,2);default:0" json:"pajak_persen"`

	PaymentType   string `gorm:"type:varchar(20);default:'PRIBADI'" json:"payment_type"`
	QrisImage     string `gorm:"type:varchar(255)" json:"qris_image"`
	QrisName      string   `gorm:"type:varchar(100)" json:"qris_name"`

	ReceiptFooter string `gorm:"type:text" json:"receipt_footer"`

	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}