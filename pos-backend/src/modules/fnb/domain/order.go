package domain

import "time"

// OrderFnB untuk menyimpan data transaksi/nota kasir
type OrderFnB struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    StoreID      uint           `json:"store_id"`      // WAJIB buat Multitenant
    CreatedBy    uint           `json:"created_by"`    // WAJIB buat Tracking User
    Invoice      string         `json:"invoice"`
    SessionID    int            `json:"session_id"`
    TipeOrder    string         `json:"tipe_order"`
    NomorMeja    string         `json:"nomor_meja"`
    NamaPemesan  string         `json:"nama_pemesan"`
    MetodeBayar  string         `json:"metode_bayar"`
    UangDiterima int            `json:"uang_diterima"`
    Kembalian    int            `json:"kembalian"`
    TotalHarga   int            `json:"total_harga"`
    StatusDapur  string         `json:"status_dapur"`
    Items        []OrderItemFnB `gorm:"foreignKey:OrderID" json:"items"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
}

// OrderItemFnB untuk rincian menu
type OrderItemFnB struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    OrderID   uint   `json:"order_id"`
    ProductID uint   `json:"product_id"`
    Qty       int    `json:"qty"`
    SubTotal  int    `json:"sub_total"`
    Notes     string `json:"notes"`
}

func (OrderFnB) TableName() string {
    return "orders_fnb" // NAMA TABEL BARU YANG BERSIH
}

func (OrderItemFnB) TableName() string {
    return "order_items_fnb" // NAMA TABEL BARU YANG BERSIH
}