package domain

import (
	"pos-backend/models"
	"time"
)

type TransactionLaundryDetail struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	StoreID       uint      `gorm:"index;not null" json:"store_id"`
	TransactionID uint      `gorm:"index;not null" json:"transaction_id"`
	ProductID     uint      `gorm:"index;not null" json:"product_id"`
	NamaPelanggan string    `gorm:"type:varchar(100);not null" json:"nama_pelanggan"`
	NoWhatsapp    string    `gorm:"type:varchar(20);not null" json:"no_whatsapp"`
	BeratKg       float64   `gorm:"type:decimal(10,2);not null" json:"berat_kg"`
	HargaPerKg    float64   `gorm:"type:decimal(10,2);not null" json:"harga_per_kg"`
	NamaParfum    string    `gorm:"type:varchar(100);default:'Parfum Standar'" json:"nama_parfum"`
	HargaParfum   float64   `gorm:"type:decimal(10,2);default:0" json:"harga_parfum"`
	IsExpress     bool      `gorm:"default:false" json:"is_express"`
	HargaExpress  float64   `gorm:"type:decimal(10,2);default:0" json:"harga_express"`
	RackID        uint      `gorm:"index" json:"rack_id"`
	NomorRak      string    `gorm:"type:varchar(50);default:'-'" json:"nomor_rak"`
	SubTotal      float64   `gorm:"type:decimal(10,2);not null" json:"sub_total"`
	StatusBayar   string    `gorm:"type:varchar(30);index;default:'BELUM_BAYAR'" json:"status_bayar"`
	MetodeBayar   string    `gorm:"type:varchar(50);default:'CASH'" json:"metode_bayar"`
	FotoBarang    string    `gorm:"type:varchar(255)" json:"foto_barang"`
	StatusCucian  string    `gorm:"type:varchar(50);index;default:'ANTRI'" json:"status_cucian"`
	EstimasiWaktu time.Time `gorm:"type:timestamp" json:"estimasi_waktu"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Product models.Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"product"`
}

type TrackingResponse struct {
	ID           uint    `json:"id"`
	Invoice      string  `json:"invoice"`
	Pelanggan    string  `json:"pelanggan"`
	Whatsapp     string  `json:"whatsapp"`
	Layanan      string  `json:"layanan"`
	BeratKg      float64 `json:"berat_kg"`
	SubTotal     float64 `json:"sub_total"`
	StatusCucian string  `json:"status_cucian"`
	StatusBayar  string  `json:"status_bayar"`

	RackID   uint   `json:"rack_id"`
	NomorRak string `json:"nomor_rak"`

	EstimasiWaktu time.Time `json:"estimasi_waktu"`
}

type LaundryItemInput struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	BeratKg     float64 `json:"berat_kg" binding:"required"`
	HargaPerKg  float64 `json:"harga_per_kg" binding:"required"`
	SubTotal    float64 `json:"sub_total" binding:"required"`
	NamaParfum  string  `json:"nama_parfum"`
	HargaParfum float64 `json:"harga_parfum"`
}

type CheckoutLaundryInput struct {
	CustomerName        string             `json:"customer_name" binding:"required"`
	CustomerPhone       string             `json:"customer_phone" binding:"required"`
	EstimasiSelesai     string             `json:"estimasi_selesai" binding:"required"`
	Items               []LaundryItemInput `json:"items" binding:"required"`
	TotalAmount         float64            `json:"total_amount" binding:"required"`
	PaymentMethod       string             `json:"payment_method" binding:"required"`
	PaymentStatus       string             `json:"payment_status" binding:"required"`
	FotoBarangBase64    string             `json:"foto_barang_base64"`
	BuktiTransferBase64 string             `json:"bukti_transfer_base64"`
}

type PelunasanInput struct {
	MetodeBayar         string `json:"metode_bayar" binding:"required"`
	BuktiTransferBase64 string `json:"bukti_transfer_base64"`
}

type ReportSummary struct {
	TotalOmset float64 `json:"total_omset"`
	TotalOrder int     `json:"total_order"`
	RataRata   float64 `json:"rata_rata"`
	Tunai      float64 `json:"tunai"`
	Qris       float64 `json:"qris"`
	Debit      float64 `json:"debit"`
	Piutang    float64 `json:"piutang"`
}

type TransactionReportResponse struct {
	models.Transaction
	Invoice       string    `json:"invoice"`
	Pelanggan     string    `json:"pelanggan"`
	Whatsapp      string    `json:"whatsapp"`
	Layanan       string    `json:"layanan"`
	BeratKg       float64   `json:"berat_kg"`
	SatuanDasar   string    `json:"satuan_dasar"`
	SubTotal      float64   `json:"sub_total"`
	EstimasiWaktu time.Time `json:"estimasi_waktu"`
	NomorRak      string    `json:"nomor_rak"`
}

type ReportSummaryResponse struct {
	Ringkasan ReportSummary               `json:"ringkasan"`
	Transaksi []TransactionReportResponse `json:"transaksi"`
}

type LaundryRack struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoreID   uint      `gorm:"index;not null" json:"store_id"`
	Zona      string    `gorm:"type:varchar(100);not null;default:'Rak Utama'" json:"zona"`
	NamaRak   string    `gorm:"type:varchar(50);not null" json:"nama_rak"`         // Cth: "A-1", "B-2"
	Baris     int       `json:"baris"`                                             // Posisi Koordinat Y
	Kolom     int       `json:"kolom"`                                             // Posisi Koordinat X
	Status    string    `gorm:"type:varchar(20);default:'TERSEDIA'" json:"status"` // TERSEDIA / RUSAK
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TransactionLaundryDetail) TableName() string {
	return "laundry_transaction_details"
}

func (LaundryRack) TableName() string {
	return "laundry_racks"
}
