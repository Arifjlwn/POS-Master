package controllers

import (
    "net/http"
    "pos-backend/config"
    "pos-backend/models"
    "time"

    "github.com/gin-gonic/gin"
)

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