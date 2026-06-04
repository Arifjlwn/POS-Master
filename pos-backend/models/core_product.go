package models

import "time"

type Product struct {
	// ==========================================
	// PRIMARY IDENTIFIER
	// ==========================================
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex;not null" json:"public_id"`

	// ==========================================
	// MULTI TENANT
	// ==========================================
	StoreID uint `gorm:"not null;index" json:"store_id"`

	// ==========================================
	// INFORMASI PRODUK
	// ==========================================
	SKU *string `gorm:"type:varchar(50);index" json:"sku"`

	NamaProduk string `gorm:"type:varchar(150);not null;index" json:"nama_produk"`

	Kategori string `gorm:"type:varchar(50);index" json:"kategori"`

	ProductType string `gorm:"type:varchar(20);default:'retail';index" json:"product_type"`

	Estimasi string `gorm:"type:varchar(50);default:'Standar'" json:"estimasi"`

	// ==========================================
	// HARGA
	// ==========================================
	HargaModal float64 `gorm:"type:decimal(18,2);default:0" json:"harga_modal"`

	HargaJual float64 `gorm:"type:decimal(18,2);not null" json:"harga_jual"`

	// ==========================================
	// STOK
	// ==========================================
	Stok int `gorm:"default:0" json:"stok"`

	IsAvailable bool `gorm:"default:true" json:"is_available"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	// ==========================================
	// MEDIA
	// ==========================================
	Gambar string `gorm:"type:text" json:"gambar"`

	// ==========================================
	// SATUAN DASAR
	// ==========================================
	SatuanDasar string `gorm:"type:varchar(20);default:'PCS'" json:"satuan_dasar"`

	// ==========================================
	// KONVERSI 2 LAPIS
	// ==========================================
	SatuanBesar string `gorm:"type:varchar(20)" json:"satuan_besar"`

	IsiPerBesar int `gorm:"default:0" json:"isi_per_besar"`

	HargaJualBesar float64 `gorm:"type:decimal(18,2);default:0" json:"harga_jual_besar"`

	// ==========================================
	// KONVERSI 3 LAPIS
	// ==========================================
	IsNestedUom bool `gorm:"default:false" json:"is_nested_uom"`

	SatuanTengah string `gorm:"type:varchar(20)" json:"satuan_tengah"`

	IsiBesarKeTengah int `gorm:"default:0" json:"isi_besar_ke_tengah"`

	IsiTengahKeDasar int `gorm:"default:0" json:"isi_tengah_ke_dasar"`

	HargaJualTengah float64 `gorm:"type:decimal(18,2);default:0" json:"harga_jual_tengah"`

	// ==========================================
	// RELATION
	// ==========================================
	Store Store `gorm:"foreignKey:StoreID" json:"-"`

	// ==========================================
	// TIMESTAMP
	// ==========================================
	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}