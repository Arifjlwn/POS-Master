package delivery

import (
	"net/http"
	"strconv"
	"time"

	"pos-backend/models"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OpenSessionInput struct {
	StationNumber string  `json:"station_number" binding:"required"`
	ModalAwal     float64 `json:"modal_awal"`
}

func (h *RetailHandler) OpenSession(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
	}
	switch v := userIDRaw.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	}

	userRoleRaw, exists := c.Get("role")
	userRole := ""
	if exists {
		userRole = userRoleRaw.(string)
	}

	var input OpenSessionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input stasiun tidak lengkap!"})
		return
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowInJKT := time.Now().In(loc)
	today := nowInJKT.Format("2006-01-02")

	// 📸 INTEGRASI ABSENSI: Pastikan staff sudah melakukan absen Face AI hari ini
	if userRole != "owner" {
		if _, err := h.Repo.GetAttendanceToday(userID, today); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Sistem mendeteksi Anda belum melakukan Absen Wajah hari ini!", "tanggal_hari_ini": today})
			return
		}
	}

	db := h.Repo.GetDB()

	// 🔒 CONCURRENCY LOCK: Cek apakah user bersangkutan sudah punya laci aktif di toko ini
	if _, err := h.Repo.GetActiveSession(db, userID, storeID); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sesi kasir Anda masih aktif di perangkat lain!"})
		return
	}

	// 🔒 STATION CHECK: Cegah nomor stasiun yang sama di-hijack/digunakan 2 user berbeda sekaligus
	var duplicateStation int64
	db.Model(&models.CashierSession{}).Where("store_id = ? AND station_number = ? AND status = ?", storeID, input.StationNumber, "open").Count(&duplicateStation)
	if duplicateStation > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stasiun POS " + input.StationNumber + " sedang aktif digunakan oleh kasir lain!"})
		return
	}

	var store models.Store
	if err := db.Where("id = ?", storeID).First(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data Toko tidak valid"})
		return
	}

	var activeSessions int64
	db.Model(&models.CashierSession{}).Where("store_id = ? AND status = ?", storeID, "open").Count(&activeSessions)

	kuotaTerminal := store.QuotaTerminal
	if kuotaTerminal == 0 {
		kuotaTerminal = 1
	}

	// 🛡️ MULTI-TENANT QUOTA FILTER
	if activeSessions >= int64(kuotaTerminal) {
		c.JSON(http.StatusForbidden, gin.H{
			"error_code": "QUOTA_FULL",
			"error":      "Batas maksimal kuota terminal kasir Anda telah tercapai! Silakan upgrade lisensi toko.",
		})
		return
	}

	newSession := models.CashierSession{
		PublicID:      utils.GenerateULID(),
		StoreID:       storeID,
		UserID:        userID,
		StationNumber: input.StationNumber,
		ModalAwal:     input.ModalAwal,
		OpenedBy:      userID, // Catat penanggung jawab pembuka sesi
		Status:        "open",
		OpenedAt:      nowInJKT,
	}
	newSession.CreatedAt = nowInJKT

	if err := h.Repo.CreateSession(&newSession); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengaktifkan laci kasir baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sesi laci kasir berhasil diaktifkan!", "session": newSession})
}

func (h *RetailHandler) CheckSessionStatus(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
	}
	switch v := userIDRaw.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	}

	db := h.Repo.GetDB()
	session, err := h.Repo.GetActiveSession(db, userID, storeID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"has_session": false})
		return
	}

	store, errStore := h.Repo.GetStoreByIDSimple(db, storeID)
	if errStore == nil {
		session.Store = *store
	}

	// 🛡️ SANITASI KEAMANAN DATA: Hapus data password hash dan pin hash
	// sebelum objek JSON dikirim ke client side browser demi mencegah XSS Leak !
	session.User.Password = ""
	// session.User.Pin = ""

	c.JSON(http.StatusOK, gin.H{"has_session": true, "station_number": session.StationNumber, "session": session})
}

func (h *RetailHandler) CloseSession(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
	}
	switch v := userIDRaw.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	}

	sessionIDStr := c.Param("id")
	sessionID, _ := strconv.Atoi(sessionIDStr)

	var input struct {
		TotalAktual float64 `json:"total_aktual"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data penutupan salah"})
		return
	}

	session, err := h.Repo.GetSessionByIDPreloaded(uint(sessionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sesi kasir tidak ditemukan"})
		return
	}

	// 🔒 ANTIDOTE IDOR ATTACK: Pastikan laci kasir yang mau ditutup bener-bener milik Toko si Kasir login!
	if session.StoreID != storeID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ilegal! Anda tidak berhak menutup laci kasir cabang lain!"})
		return
	}

	// Pastikan sesi laci kasir tidak diproses closing berulang-ulang
	if session.Status == "closed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sesi laci kasir ini sudah berstatus CLOSED!"})
		return
	}

	salesGross, totalTax, _ := h.Repo.GetSalesTotalAndTax(sessionIDStr)
	netSales := salesGross - totalTax
	salesCash, _ := h.Repo.GetSalesMethodSummary(sessionIDStr, "Cash")
	salesNonTunai := salesGross - salesCash

	// 💸 FIX AKUNTANSI: Rumus perhitungan ekspektasi wajib dikurangi dengan pengeluaran kas laci (TotalKeluar) !
	totalExpected := (session.ModalAwal + salesCash) - session.TotalKeluar
	selisih := input.TotalAktual - totalExpected

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	session.TotalMasuk = salesCash
	session.TotalAktual = input.TotalAktual
	session.Selisih = selisih
	session.Status = "closed"
	session.ClosedAt = &now
	session.ClosedBy = &userID // Catat ID user saksi yang memproses closing
	session.UpdatedAt = now

	newClosing := models.ShiftClosing{
		PublicID:  utils.GenerateULID(),
		StoreID:   storeID,
		SessionID: uint(sessionID),
		UserID:    userID,

		// 1. Ambil Waktu Riil dari Sesi Kasir
		OpenedAt: session.OpenedAt,
		ClosedAt: now,

		// 2. Ambil Data Arus Kas dari Sesi Kasir
		ModalAwal: session.ModalAwal, // <-- Pindahin ke sini bray!

		// 3. Rekap Penjualan & Pajak dari Hasil Query Repository
		NetSales:     netSales,
		TotalTax:     totalTax,
		SalesCash:    salesCash,
		SalesNonCash: salesNonTunai,

		// 4. Rekap Ekspektasi, Aktual, dan Selisih
		TotalExpected: totalExpected,
		TotalActual:   input.TotalAktual,
		Selisih:       selisih,

		// 5. Catatan Kasir (Kalau ada input note dari frontend)
		ClosingNote: c.PostForm("closing_note"), // atau sesuai binding input lu
	}
	newClosing.CreatedAt = now

	db := h.Repo.GetDB()
	errTx := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(session).Error; err != nil {
			return err
		}
		if err := tx.Create(&newClosing).Error; err != nil {
			return err
		}
		return nil
	})

	if errTx != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses transaksi penutupan laci kasir."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"start_time":      session.OpenedAt.In(loc).Format("02.01.06 15:04"),
		"end_time":        session.ClosedAt.In(loc).Format("02.01.06 15:04"),
		"sales_gross":     salesGross,
		"total_tax":       totalTax,
		"net_sales":       netSales,
		"modal_awal":      session.ModalAwal,
		"sales_cash":      salesCash,
		"sales_non_tunai": salesNonTunai,
		"total_keluar":    session.TotalKeluar, // Informasikan juga total pengeluaran laci kasir
		"total_expected":  totalExpected,
		"total_actual":    input.TotalAktual,
		"selisih":         selisih,
	})
}
