package models

import (
	"time"

	"github.com/oklog/ulid/v2" // Pastikan lu install ini atau sesuaikan dengan generator ULID lu
	"gorm.io/gorm"
)

type Product struct {
	// ==========================================
	// PRIMARY IDENTIFIER
	// ==========================================
	ID       uint   `gorm:"primaryKey" json:"id"`
	PublicID string `gorm:"size:26;uniqueIndex;not null" json:"public_id"`

	// ==========================================
	// MULTI TENANT (COMPOSITE INDEX WITH SKU)
	// ==========================================
	// FIX: StoreID diikat bareng SKU agar barcode unik per tenant!
	StoreID uint `gorm:"not null;uniqueIndex:idx_store_sku;index" json:"store_id"`

	// ==========================================
	// INFORMASI PRODUK
	// ==========================================
	// FIX: SKU menggunakan composite unique index agar anti-duplikat di dalam toko yang sama
	SKU *string `gorm:"type:varchar(50);uniqueIndex:idx_store_sku" json:"sku"`

	NamaProduk string `gorm:"type:varchar(150);not null;index" json:"nama_produk"`
	Kategori   string `gorm:"type:varchar(50);index" json:"kategori"`

	// FIX: Menyamakan nama field struct dengan penulisan json standard
	ProductType string `gorm:"type:varchar(20);default:'retail';index" json:"product_type"`
	Estimasi    string `gorm:"type:varchar(50);default:'Standar'" json:"estimasi"`

	// ==========================================
	// HARGA & NOMINAL FINANCIAL
	// ==========================================
	HargaModal float64 `gorm:"type:decimal(18,2);default:0" json:"harga_modal"`
	HargaJual  float64 `gorm:"type:decimal(18,2);not null" json:"harga_jual"`

	// ==========================================
	// STOK & STATUS RUNTIME
	// ==========================================
	Stok        int  `gorm:"default:0" json:"stok"`
	IsAvailable bool `gorm:"default:true" json:"is_available"`
	IsActive    bool `gorm:"default:true" json:"is_active"`

	// ==========================================
	// MEDIA ASSETS
	// ==========================================
	Gambar string `gorm:"type:text" json:"gambar"`

	// ==========================================
	// SATUAN DASAR (UOM LAYER 1)
	// ==========================================
	SatuanDasar string `gorm:"type:varchar(20);default:'PCS'" json:"satuan_dasar"`
	EstimasiDurasi int     `gorm:"default:1" json:"estimasi_durasi"`
	EstimasiSatuan string  `gorm:"size:10;default:'Hari'" json:"estimasi_satuan"`

	// ==========================================
	// KONVERSI 2 LAPIS (UOM LAYER 2)
	// ==========================================
	SatuanBesar    string  `gorm:"type:varchar(20)" json:"satuan_besar"`
	IsiPerBesar    int     `gorm:"default:0" json:"isi_per_besar"`
	HargaJualBesar float64 `gorm:"type:decimal(18,2);default:0" json:"harga_jual_besar"`

	// ==========================================
	// KONVERSI 3 LAPIS (UOM LAYER 3)
	// ==========================================
	IsNestedUom      bool    `gorm:"default:false" json:"is_nested_uom"`
	SatuanTengah     string  `gorm:"type:varchar(20)" json:"satuan_tengah"`
	IsiBesarKeTengah int     `gorm:"default:0" json:"isi_besar_ke_tengah"`
	IsiTengahKeDasar int     `gorm:"default:0" json:"isi_tengah_ke_dasar"`
	HargaJualTengah  float64 `gorm:"type:decimal(18,2);default:0" json:"harga_jual_tengah"`

	// ==========================================
	// RELATIONSHIPS
	// ==========================================
	Store Store `gorm:"foreignKey:StoreID" json:"-"`

	// ==========================================
	// TIMESTAMPS
	// ==========================================
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 🚀 BONUS SAAS ENTERPRISE: GORM Hook untuk auto-generate ULID Public ID !
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.PublicID == "" {
		p.PublicID = ulid.Make().String()
	}
	return
}
