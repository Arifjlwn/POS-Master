package config

import (
	"log"
	"os"
	"pos-backend/models" // Pastikan path ini benar sesuai struktur folder kamu

	"gorm.io/driver/postgres" 
	"gorm.io/gorm"
)

// Variabel Global untuk menampung Koneksi DB
var DB *gorm.DB

func ConnectDatabase() {
	// Membaca string koneksi URI dari file .env
	dsn := os.Getenv("DB_URL")
	
	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Gagal Menyambung ke Supabase! Error: ", err)
	}

	log.Println("✅ Berhasil terhubung ke Database Supabase (Postgres)!")

	// Auto Migrate (GORM bakal otomatis ngebangun/update tabel di Cloud Supabase)
	err = database.AutoMigrate(
		
		// ==========================================
		// 🌐 CORE GLOBAL (Pondasi Utama & SDM)
		// Dipakai oleh semua tipe bisnis
		// ==========================================
		&models.Store{},
		&models.User{},
		&models.Product{},        // Master Katalog (Barang Fisik / Layanan Jasa)
		&models.Transaction{},    // Induk Transaksi
		&models.CashierSession{}, // Sesi Buka/Tutup Kasir
		&models.Attendance{},     // Absensi Pegawai
		&models.Schedule{},       // Jadwal Shift
		&models.Customer{},		  // CRM (Customer Relationship Management)
		
		// ==========================================
		// 🛍️ RETAIL & INVENTORY
		// Khusus manajemen barang fisik berskala gudang
		// ==========================================
		&models.TransactionDetail{}, // Anak transaksi khusus Retail
		&models.Purchase{},          // Pembelian (LPB)
		&models.PurchaseDetail{},
		&models.StockOpname{},       // Cek Fisik Gudang
		&models.StockOpnameDetail{},
		&models.ProductReturn{},     // Retur Barang
		
		// ==========================================
		// 🧺 LAYANAN & JASA (LAUNDRY)
		// Khusus transaksi menggunakan berat (Kg) & Status
		// ==========================================
		&models.TransactionLaundryDetail{},
		

	)
	
	if err != nil {
		log.Fatal("Gagal Melakukan Migrasi Database! Error: ", err)
	}
	log.Println("✅ Tabel database berhasil di-generate di Cloud Supabase!")

	DB = database
}