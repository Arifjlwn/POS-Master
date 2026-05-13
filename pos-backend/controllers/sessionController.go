package controllers

import (
    "net/http"
    "pos-backend/config"
    "pos-backend/models"
    "time"

    "github.com/gin-gonic/gin"
)

// OPEN SESSION
type OpenSessionInput struct {
    StationNumber string  `json:"station_number" binding:"required"`
    ModalAwal     float64 `json:"modal_awal"`
}

func OpenSession(c *gin.Context) {
    userIDRaw, _ := c.Get("user_id")
    userID := uint(userIDRaw.(float64))
    storeIDRaw, _ := c.Get("store_id")
    storeID := uint(storeIDRaw.(float64))

    var input OpenSessionInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap!"})
        return
    }

    // 🚀 1. SETTING TIMEZONE JAKARTA (WIB)
    loc, _ := time.LoadLocation("Asia/Jakarta")
    nowInJKT := time.Now().In(loc) 
    today := nowInJKT.Format("2006-01-02")

    // 2. VALIDASI: Apakah sudah absen masuk hari ini?
    var attendance models.Attendance
    if err := config.DB.Where("user_id = ? AND tanggal = ?", userID, today).First(&attendance).Error; err != nil {
        // Kita beri tahu user tanggal server vs tanggal browser
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Anda wajib Absen Wajah terlebih dahulu!",
            "tanggal_hari_ini": today,
        })
        return
    }

    // 3. VALIDASI: Apakah ada session yang masih menggantung (Open)?
    var existingSession models.CashierSession
    if err := config.DB.Where("user_id = ? AND status = ?", userID, "open").First(&existingSession).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Anda masih memiliki session yang terbuka!",
            "session_id": existingSession.ID,
        })
        return
    }

    // 4. EKSEKUSI: Buat Session Baru
    newSession := models.CashierSession{
        StoreID:       storeID,
        UserID:        userID,
        StationNumber: input.StationNumber,
        ModalAwal:     input.ModalAwal,
        StartTime:     nowInJKT, // 👈 Pakai waktu JKT biar sinkron dengan absen
        Status:        "open",
    }

    if err := config.DB.Create(&newSession).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka session kasir"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Kasir berhasil dibuka! Selamat bertugas.",
        "session": newSession,
    })
}

// CHECK SESSION
func CheckSessionStatus(c *gin.Context) {
    userIDRaw, _ := c.Get("user_id")
    userID := uint(userIDRaw.(float64))

    var session models.CashierSession
    // Kita ambil session terakhir yang statusnya open
    if err := config.DB.Where("user_id = ? AND status = ?", userID, "open").Order("id desc").First(&session).Error; err != nil {
        c.JSON(http.StatusOK, gin.H{"has_session": false})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "has_session": true,
        "session":     session,
    })
}

// CLOSE SESSION
type CloseSessionInput struct {
	TotalAktual float64 `json:"total_aktual"`
}

// 🏁 2. FUNGSI CLOSING SHIFT (TUTUP KASIR)
func CloseSession(c *gin.Context) {
	sessionID := c.Param("id")
	var input struct {
		TotalAktual float64 `json:"total_aktual"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah"})
		return
	}

	// 1. Cari Session
	var session models.CashierSession
	if err := config.DB.Preload("Store").First(&session, sessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session tidak ditemukan"})
		return
	}

	// 🚀 2. HITUNG SALES DARI DATABASE (Real-time)
	var salesGross, totalTax, salesCash, salesQRIS float64

	// Ambil semua transaksi yang terjadi dalam session ini
	// Asumsi: Mas punya tabel 'transactions' dengan kolom 'session_id', 'total_harga', 'pajak', dan 'metode_bayar'
	type SalesSummary struct {
		Gross float64
		Tax   float64
	}
	var summary SalesSummary
	
	// 1. Hitung Total Penjualan & Pajak
	config.DB.Table("transactions").
		Select("SUM(total_harga) as gross, SUM(pajak) as tax").
		Where("session_id = ?", sessionID).
		Scan(&summary)

	salesGross = summary.Gross
	totalTax = summary.Tax
	netSales := salesGross - totalTax

	// 1. Hitung khusus pembayaran TUNAI (CASH) untuk laci
    config.DB.Table("transactions").
        Select("COALESCE(SUM(total_harga), 0)").
        Where("session_id = ? AND metode_bayar = ?", sessionID, "Cash").
        Scan(&salesCash)

    // 🚀 2. Hitung Sales QRIS secara Spesifik
    config.DB.Table("transactions").
        Select("COALESCE(SUM(total_harga), 0)").
        Where("session_id = ? AND metode_bayar = ?", sessionID, "QRIS").
        Scan(&salesQRIS)

    // 🚀 3. Hitung Sales BANK / DEBIT (Contoh: BCA)
    var salesBCA float64
    config.DB.Table("transactions").
        Select("COALESCE(SUM(total_harga), 0)").
        Where("session_id = ? AND metode_bayar = ?", sessionID, "BCA").
        Scan(&salesBCA)

    // Total Non-Tunai (Gabungan semua selain Cash)
    salesNonTunai := salesQRIS + salesBCA

	// 4. Logic Akuntansi
	totalExpected := session.ModalAwal + salesCash
	selisih := input.TotalAktual - totalExpected

	// 🕒 Waktu Jakarta
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	// 5. Update Database
	session.TotalMasuk = salesCash // Uang tunai yang masuk
	session.TotalAktual = input.TotalAktual
	session.Selisih = selisih
	session.EndTime = &now
	session.Status = "closed"

	config.DB.Save(&session)

	// 5. Kirim respon (Tambahkan sales_qris ke JSON biar Vue bisa nampilin)
    c.JSON(http.StatusOK, gin.H{
        "start_time":     session.StartTime.In(loc).Format("02.01.06 15:04"),
        "end_time":       session.EndTime.In(loc).Format("02.01.06 15:04"),
        "sales_gross":    salesGross,
        "total_tax":      totalTax,
        "net_sales":      netSales,
        "modal_awal":     session.ModalAwal,
        "sales_cash":       salesCash,
        "sales_qris":       salesQRIS,
        "sales_bca":        salesBCA,
        "sales_non_tunai":  salesNonTunai,
        "total_expected":   session.ModalAwal + salesCash,
        "total_actual":   input.TotalAktual,
        "selisih":        selisih,
    })
}

