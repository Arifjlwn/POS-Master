package src

import (
    "log"
    "os"
    
    retailDomain "pos-backend/src/modules/retail/domain"
    fnbDomain"pos-backend/src/modules/fnb/domain"
    laundryDomain "pos-backend/src/modules/jasalayanan/laundry/domain"
    "pos-backend/models" // Biarin aja, buat nahan modul Retail & Laundry biar ga error
    "gorm.io/driver/postgres" 
    "gorm.io/gorm"
)

// Variabel Global untuk menampung Koneksi DB
var DB *gorm.DB

func ConnectDatabase() {
    dsn := os.Getenv("DB_URL")
    
    database, err := gorm.Open(postgres.New(postgres.Config{
        DSN:                  dsn,
        PreferSimpleProtocol: true,
    }), &gorm.Config{})

    if err != nil {
        log.Fatal("Gagal Menyambung ke Supabase! Error: ", err)
    }

    log.Println("✅ Berhasil terhubung ke Database Supabase (Postgres)!")

    err = database.AutoMigrate(
        // ==========================================
        // 🌐 CORE GLOBAL (Pondasi Utama & SDM)
        // ==========================================
        &models.Store{},
        &models.User{},
        &models.Product{},        
        &models.Transaction{},    
        &models.CashierSession{}, 
        &models.Attendance{},     
        &models.Schedule{},       
        &models.Customer{},       
        
        // ==========================================
        // 🛍️ RETAIL & INVENTORY
        // ==========================================
        &models.TransactionDetail{}, 
        &retailDomain.StockOpname{},
        &retailDomain.StockOpnameDetail{},
        &retailDomain.ProductReturn{},
        &retailDomain.Purchase{},
        &retailDomain.PurchaseDetail{},     
        
        // ==========================================
        // 🧺 LAYANAN & JASA (LAUNDRY)
        // ==========================================
        &laundryDomain.TransactionLaundryDetail{},
        &laundryDomain.Perfume{},

        // ==========================================
        // 🍔 FOOD & BEVERAGES (MODULAR BARU)
        // ==========================================
        &fnbDomain.Menu{},         // 👈 INI YANG TADI KETINGGALAN BEB!
        &fnbDomain.OrderFnB{},
        &fnbDomain.OrderItemFnB{},
    )
    
    if err != nil {
        log.Fatal("Gagal Melakukan Migrasi Database! Error: ", err)
    }
    log.Println("✅ Tabel database berhasil di-generate di Cloud Supabase!")

    DB = database
}