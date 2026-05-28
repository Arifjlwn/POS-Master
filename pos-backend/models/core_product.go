package models

import "time"

type Product struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    StoreID      uint      `gorm:"not null" json:"store_id"` // Kunci agar produk nggak nyasar ke toko lain
    SKU          *string   `gorm:"type:varchar(50);uniqueIndex" json:"sku"`  // Barcode/Kode Barang (Pakai * biar bisa NULL buat Cafe/Jasa)
    NamaProduk   string    `gorm:"type:varchar(150);not null" json:"nama_produk"`
    Kategori     string    `gorm:"type:varchar(50)" json:"kategori"`
    Estimasi     string    `gorm:"type:varchar(50);default:'Standar'" json:"estimasi"`
    HargaModal   float64   `gorm:"type:decimal(10,2);default:0" json:"harga_modal"`
    HargaJual    float64   `gorm:"type:decimal(10,2);not null" json:"harga_jual"`
    Stok         int       `gorm:"default:0" json:"stok"`
    Gambar       string    `gorm:"type:varchar(255)" json:"gambar"`
    IsAvailable  bool      `gorm:"default:true" json:"is_available"` // True = Menu Ada, False = Habis/Sembunyikan
    
    // 🚀 PASUKAN KONVERSI 2 LAPIS (KARTON/KARUNG)
    SatuanDasar    string   `gorm:"type:varchar(20);default:'PCS'" json:"satuan_dasar"`
    SatuanBesar    string   `gorm:"type:varchar(20)" json:"satuan_besar"`
    IsiPerBesar    int      `gorm:"default:0" json:"isi_per_besar"`
    HargaJualBesar float64  `gorm:"type:decimal(10,2);default:0" json:"harga_jual_besar"`

    // 🚀 INI DIA AMUNISI SILUMAN BARU: KONVERSI 3 LAPIS (ROKOK/RENTENG)
    IsNestedUom      bool    `gorm:"default:false" json:"is_nested_uom"`
    SatuanTengah     string  `gorm:"type:varchar(20)" json:"satuan_tengah"`
    IsiBesarKeTengah int     `gorm:"default:0" json:"isi_besar_ke_tengah"`
    IsiTengahKeDasar int     `gorm:"default:0" json:"isi_tengah_ke_dasar"`
    HargaJualTengah  float64 `gorm:"type:decimal(10,2);default:0" json:"harga_jual_tengah"`

    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`

    // Relasi ke Store
    Store Store `gorm:"foreignKey:StoreID" json:"-"`
}