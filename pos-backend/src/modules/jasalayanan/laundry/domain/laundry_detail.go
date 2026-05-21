package domain

import (
	"time"
	"pos-backend/models" // Relasi Product global
)

type TransactionLaundryDetail struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	StoreID       uint      `gorm:"index;not null" json:"store_id"`
	TransactionID uint      `gorm:"index" json:"transaction_id"`   
	ProductID     uint      `json:"product_id"`                     
	NamaPelanggan string    `gorm:"type:varchar(100)" json:"nama_pelanggan"`
	NoWhatsapp    string    `gorm:"type:varchar(20)" json:"no_whatsapp"`
	BeratKg       float64   `json:"berat_kg"`       
	HargaPerKg    float64   `json:"harga_per_kg"`   
	SubTotal      float64   `json:"sub_total"`
	FotoBarang    string    `gorm:"type:varchar(255)" json:"foto_barang"`
	StatusCucian  string    `gorm:"type:varchar(50);index;default:'ANTRI'" json:"status_cucian"`
	EstimasiWaktu time.Time `json:"estimasi_waktu"` 
	NamaParfum    string    `gorm:"type:varchar(100);default:'Parfum Standar'" json:"nama_parfum"`
	HargaParfum   float64   `gorm:"type:decimal(10,2);default:0" json:"harga_parfum"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Product       models.Product `gorm:"foreignKey:ProductID" json:"product"`
}

// Wadah kustom penampung hasil Join 3 Tabel (Ganti nama tabel sesuai domain baru)
type TrackingResponse struct {
	ID        uint    `json:"id"`
	Invoice   string  `json:"invoice"`
	Pelanggan string  `json:"pelanggan"`
	Whatsapp  string  `json:"whatsapp"`
	Layanan   string  `json:"layanan"`
	BeratKg   float64 `json:"berat_kg"`
	SubTotal  float64 `json:"sub_total"`
	Status    string  `json:"status"`
}

func (TransactionLaundryDetail) TableName() string {
	return "laundry_transaction_details"
}