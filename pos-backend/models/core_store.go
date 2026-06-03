package models

import (
	"time"
)

// Store merepresentasikan tabel 'stores' di Supabase yang sinkron dengan Form Setup
type Store struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	OwnerID  uint   `gorm:"index" json:"owner_id"` // TAMBAHAN PENTING UNTUK MULTI-OUTLET
	NamaToko string `gorm:"type:varchar(100);not null" json:"nama_toko"`
	Telepon  string `gorm:"type:varchar(20)" json:"telepon"`

	BusinessType string `gorm:"type:varchar(100);column:business_type;not null" json:"business_type"`
	Industry     string `gorm:"type:varchar(50);default:'retail'" json:"industry"`

	SubscriptionPlan   string    `gorm:"type:varchar(50);default:'trial'" json:"subscription_plan"`
	SubscriptionStatus string    `gorm:"type:varchar(20);default:'active'" json:"subscription_status"`
	SubscriptionEnd    time.Time `json:"subscription_end"`

	FiturAktif    string `gorm:"type:text;column:fitur_aktif" json:"fitur_aktif"`
	QuotaTerminal int    `json:"quota_terminal" gorm:"default:1"`

	// Lokasi Toko
	Alamat    string `gorm:"type:text" json:"alamat"`
	Provinsi  string `gorm:"type:varchar(100)" json:"provinsi"`
	Kota      string `gorm:"type:varchar(100)" json:"kota"`
	Kecamatan string `gorm:"type:varchar(100)" json:"kecamatan"`
	Kelurahan string `gorm:"type:varchar(100)" json:"kelurahan"`
	KodePos   string `gorm:"type:varchar(10)" json:"kode_pos"`

	// Branding Toko & Pajak
	LogoURL       string  `gorm:"type:text" json:"logo_url"`
	IsTaxActive   bool    `gorm:"default:false" json:"is_tax_active"`
	PajakPersen   float64 `gorm:"type:decimal(5,2);default:0" json:"pajak_persen"`
	ReceiptFooter string  `gorm:"type:text" json:"receipt_footer"`

	// PENGATURAN PEMBAYARAN (QRIS & MIDTRANS)
	PaymentType       string `gorm:"type:varchar(50);default:'qris_static'" json:"payment_type"`
	QrisImage         string `gorm:"type:varchar(255)" json:"qris_image"`
	QrisName          string `gorm:"type:varchar(100)" json:"qris_name"`
	MidtransServerKey string `gorm:"type:varchar(255)" json:"midtrans_server_key"`
	MidtransClientKey string `gorm:"type:varchar(255)" json:"midtrans_client_key"`

	// PENGATURAN PRINTER THERMAL
	PrinterWidth string `gorm:"type:varchar(20);default:'58mm'" json:"printer_width"`
	PrinterType  string `gorm:"type:varchar(20);default:'bluetooth'" json:"printer_type"`

	// TAMBAHAN SISTEM WA SAAS
	WaToken string `gorm:"type:varchar(255)" json:"wa_token"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
