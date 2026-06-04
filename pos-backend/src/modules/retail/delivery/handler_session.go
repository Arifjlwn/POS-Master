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

// ================================
// 💵 CASHIER SESSION HANDLERS
// ================================

type OpenSessionInput struct {
	StationNumber string  `json:"station_number" binding:"required"`
	ModalAwal     float64 `json:"modal_awal"`
}

func (h *RetailHandler) OpenSession(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id"); userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }
	switch v := userIDRaw.(type) { case float64: userID = uint(v); case uint: userID = v; case int: userID = uint(v) }

	userRoleRaw, exists := c.Get("role")
	userRole := ""
	if exists { userRole = userRoleRaw.(string) }

	var input OpenSessionInput
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap!"}); return }

	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowInJKT := time.Now().In(loc)
	today := nowInJKT.Format("2006-01-02")

	if userRole != "owner" {
		if _, err := h.Repo.GetAttendanceToday(userID, today); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda wajib Absen Wajah terlebih dahulu!", "tanggal_hari_ini": today})
			return
		}
	}

	db := h.Repo.GetDB()
	if _, err := h.Repo.GetActiveSession(db, userID, storeID); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda masih memiliki session yang terbuka!"})
		return
	}

	var store models.Store
	if err := db.Where("id = ?", storeID).First(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Toko tidak ditemukan"})
		return
	}

	var activeSessions int64
	db.Model(&models.CashierSession{}).Where("store_id = ? AND status = ?", storeID, "open").Count(&activeSessions)

	kuotaTerminal := store.QuotaTerminal
	if kuotaTerminal == 0 { kuotaTerminal = 1 }

	if activeSessions >= int64(kuotaTerminal) {
		c.JSON(http.StatusForbidden, gin.H{
			"error_code": "QUOTA_FULL",
			"error":      "Batas Terminal Kasir Tercapai! Silakan tutup shift sebelumnya atau Beli Lisensi Tambahan.",
		})
		return
	}

	newSession := models.CashierSession{
		PublicID:      utils.GenerateULID(),
		StoreID:       storeID,
		UserID:        userID,
		StationNumber: input.StationNumber,
		ModalAwal:     input.ModalAwal,
		Status:        "open",
	}
	newSession.CreatedAt = nowInJKT // 🚀 FIX: Menggunakan field bawaan GORM model (CreatedAt) buat gantiin StartTime

	if err := h.Repo.CreateSession(&newSession); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka session kasir"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kasir berhasil dibuka! Selamat bertugas.", "session": newSession})
}

func (h *RetailHandler) CheckSessionStatus(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id"); userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }
	switch v := userIDRaw.(type) { case float64: userID = uint(v); case uint: userID = v; case int: userID = uint(v) }

	db := h.Repo.GetDB()
	session, err := h.Repo.GetActiveSession(db, userID, storeID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"has_session": false})
		return
	}

	store, errStore := h.Repo.GetStoreByIDSimple(db, storeID)
	if errStore == nil { session.Store = *store }

	c.JSON(http.StatusOK, gin.H{"has_session": true, "session": session})
}

func (h *RetailHandler) CloseSession(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id"); userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }
	switch v := userIDRaw.(type) { case float64: userID = uint(v); case uint: userID = v; case int: userID = uint(v) }

	sessionIDStr := c.Param("id")
	sessionID, _ := strconv.Atoi(sessionIDStr)

	var input struct { TotalAktual float64 `json:"total_aktual"` }
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah"}); return }

	session, err := h.Repo.GetSessionByIDPreloaded(uint(sessionID))
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Session tidak ditemukan"}); return }

	salesGross, totalTax, _ := h.Repo.GetSalesTotalAndTax(sessionIDStr)
	netSales := salesGross - totalTax
	salesCash, _ := h.Repo.GetSalesMethodSummary(sessionIDStr, "Cash")
	salesNonTunai := salesGross - salesCash
	totalExpected := session.ModalAwal + salesCash
	selisih := input.TotalAktual - totalExpected

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	session.TotalMasuk = salesCash
	session.TotalAktual = input.TotalAktual
	session.Selisih = selisih
	session.UpdatedAt = now // 🚀 FIX: Menggunakan UpdatedAt bawaan GORM gantiin EndTime pointer
	session.Status = "closed"

	newClosing := models.ShiftClosing{
		PublicID:      utils.GenerateULID(),
		StoreID:       storeID,
		SessionID:     uint(sessionID),
		UserID:        userID,
		NetSales:      netSales,
		TotalTax:      totalTax,
		SalesCash:     salesCash,
		SalesNonCash:  salesNonTunai, // 🚀 FIX: Menyesuaikan penamaan standard model (SalesNonCash)
		TotalExpected: totalExpected,
		TotalActual:   input.TotalAktual,
		Selisih:       selisih,
	}
	newClosing.CreatedAt = now // 🚀 FIX: Inject waktu closing ke CreatedAt model

	db := h.Repo.GetDB()
	errTx := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(session).Error; err != nil { return err }
		if err := tx.Create(&newClosing).Error; err != nil { return err }
		return nil
	})

	if errTx != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses penutupan laci kasir secara aman."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"start_time":      session.CreatedAt.In(loc).Format("02.01.06 15:04"), // 🚀 FIX: Tarik dari CreatedAt
		"end_time":        session.UpdatedAt.In(loc).Format("02.01.06 15:04"), // 🚀 FIX: Tarik dari UpdatedAt
		"sales_gross":     salesGross,
		"total_tax":       totalTax,
		"net_sales":       netSales,
		"modal_awal":      session.ModalAwal,
		"sales_cash":      salesCash,
		"sales_non_tunai": salesNonTunai,
		"total_expected":  totalExpected,
		"total_actual":    input.TotalAktual,
		"selisih":         selisih,
	})
}