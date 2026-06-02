package delivery

import (
	"net/http"
	"strconv"
	"time"

	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// ================================
// 💵 CASHIER SESSION HANDLERS
// ================================
type OpenSessionInput struct {
	StationNumber string  `json:"station_number" binding:"required"`
	ModalAwal     float64 `json:"modal_awal"`
}

func (h *RetailHandler) OpenSession(c *gin.Context) {
	userID := uint(c.MustGet("user_id").(float64))
	storeID := uint(c.MustGet("store_id").(float64))

	userRoleRaw, exists := c.Get("role")
	userRole := ""
	if exists {
		userRole = userRoleRaw.(string)
	}

	var input OpenSessionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap!"})
		return
	}

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

	newSession := models.CashierSession{
		StoreID:       storeID,
		UserID:        userID,
		StationNumber: input.StationNumber,
		ModalAwal:     input.ModalAwal,
		StartTime:     nowInJKT,
		Status:        "open",
	}

	if err := h.Repo.CreateSession(&newSession); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka session kasir"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kasir berhasil dibuka! Selamat bertugas.", "session": newSession})
}

func (h *RetailHandler) CheckSessionStatus(c *gin.Context) {
	userID := uint(c.MustGet("user_id").(float64))
	storeID := uint(c.MustGet("store_id").(float64))

	db := h.Repo.GetDB()
	session, err := h.Repo.GetActiveSession(db, userID, storeID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"has_session": false})
		return
	}

	store, errStore := h.Repo.GetStoreByIDSimple(db, storeID)
	if errStore == nil {
		session.Store = *store // Inject data toko ke object session sebelum dikirim ke Vue
	}

	c.JSON(http.StatusOK, gin.H{"has_session": true, "session": session})
}

func (h *RetailHandler) CloseSession(c *gin.Context) {
    userID := uint(c.MustGet("user_id").(float64))     // 🚀 AMBIL USER ID
    storeID := uint(c.MustGet("store_id").(float64))   // 🚀 AMBIL STORE ID
    
    sessionIDStr := c.Param("id")
    sessionID, _ := strconv.Atoi(sessionIDStr)

    var input struct {
        TotalAktual float64 `json:"total_aktual"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah"})
        return
    }

    session, err := h.Repo.GetSessionByIDPreloaded(uint(sessionID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Session tidak ditemukan"})
        return
    }

    // Kalkulasi Total
    salesGross, totalTax, _ := h.Repo.GetSalesTotalAndTax(sessionIDStr)
    netSales := salesGross - totalTax

    salesCash, _ := h.Repo.GetSalesMethodSummary(sessionIDStr, "Cash")
    salesNonTunai := salesGross - salesCash

    totalExpected := session.ModalAwal + salesCash
    selisih := input.TotalAktual - totalExpected

    loc, _ := time.LoadLocation("Asia/Jakarta")
    now := time.Now().In(loc)

    // 1. UPDATE TABEL SESSION (Yang Lama)
    session.TotalMasuk = salesCash
    session.TotalAktual = input.TotalAktual
    session.Selisih = selisih
    session.EndTime = &now
    session.Status = "closed"
    
    // Asumsi lu punya fungsi SaveSession di repo
    h.Repo.SaveSession(session)

    // 🚀 2. CREATE KE TABEL SHIFT_CLOSING (YANG BARU!) 🚀
    newClosing := models.ShiftClosing{
        StoreID:       storeID,
        SessionID:     uint(sessionID),
        UserID:        userID,
        StartTime:     session.StartTime,
        EndTime:       now,
        NetSales:      netSales,
        TotalTax:      totalTax,
        SalesCash:     salesCash,
        SalesNonTunai: salesNonTunai,
        TotalExpected: totalExpected,
        TotalActual:   input.TotalAktual,
        Selisih:       selisih,
        // Relasi nggak usah diisi di sini, cukup ID-nya aja
    }

    // 🚀 SIMPAN KE DATABASE PAKE REPO LU
    db := h.Repo.GetDB() // Ambil instance DB lu
    if err := db.Create(&newClosing).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan laporan closing ke histori"})
        return
    }

    // 3. BALIKIN JSON KE VUE (Pastiin sesuai sama respon yang ditangkap Vue lu)
    c.JSON(http.StatusOK, gin.H{
        "start_time":      session.StartTime.In(loc).Format("02.01.06 15:04"),
        "end_time":        session.EndTime.In(loc).Format("02.01.06 15:04"),
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