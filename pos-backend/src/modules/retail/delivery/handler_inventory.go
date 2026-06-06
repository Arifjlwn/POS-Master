package delivery

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"pos-backend/models"
	"pos-backend/src/modules/retail/domain"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ==========================================
// 🚀 STOCK OPNAME (AUDIT STOK) HANDLERS
// ==========================================

type SOItemInput struct {
	ProductID uint   `json:"product_id" binding:"required"`
	ActualQty int    `json:"actual_qty" binding:"min=0"` 
	Alasan    string `json:"alasan"`
}

type StockOpnameInput struct {
	Notes string        `json:"notes"`
	Items []SOItemInput `json:"items" binding:"required,min=1"`
}

type KlaimItemInput struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Qty       int    `json:"qty" binding:"required,gt=0"`
	Alasan    string `json:"alasan" binding:"required"`
}

type KlaimRequest struct {
	Notes string           `json:"notes"`
	Items []KlaimItemInput `json:"items" binding:"required,min=1"`
}

func (h *RetailHandler) CreateStockOpname(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	userIDRaw, _ := c.Get("user_id")
	userRoleRaw, _ := c.Get("role")
	
	var storeID, userID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }
	switch v := userIDRaw.(type) { case float64: userID = uint(v); case uint: userID = v; case int: userID = uint(v) }
	userRole := strings.ToLower(userRoleRaw.(string))

	var input StockOpnameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data audit tidak valid !"})
		return
	}

	// 🌐 FIX TIMEZONE: Pakai Asia/Jakarta biar sinkron sama kasir toko bray
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	
	if h.Repo.CheckStockOpnameThisMonth(storeID, int(now.Month()), now.Year()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sistem Terkunci! Audit Stock Opname hanya bisa dilakukan 1x dalam bulan berjalan."})
		return
	}

	status := "PENDING_APPROVAL"
	if userRole == "owner" {
		status = "APPROVED"
	}

	db := h.Repo.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		so := domain.StockOpname{
			PublicID:  utils.GenerateULID(), 
			StoreID:   storeID,
			UserID:    userID,
			Notes:     input.Notes,
			Status:    status,
			CreatedAt: now,
		}
		if err := tx.Create(&so).Error; err != nil {
			return err
		}

		for _, item := range input.Items {
			var product models.Product
			
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, "id = ? AND store_id = ?", item.ProductID, storeID).Error; err != nil {
				return fmt.Errorf("produk ID %d gagal disisir di database master", item.ProductID)
			}

			calculatedSystemQty := product.Stok
			calculatedSelisih := item.ActualQty - calculatedSystemQty
			nilaiKalkulasiUang := float64(calculatedSelisih) * product.HargaModal

			detail := domain.StockOpnameDetail{
				OpnameID:  so.ID, 
				ProductID: item.ProductID,
				SystemQty: calculatedSystemQty, 
				ActualQty: item.ActualQty,
				Selisih:   calculatedSelisih,   
				NilaiUang: nilaiKalkulasiUang,
				Alasan:    item.Alasan,
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			if status == "APPROVED" {
				product.Stok = item.ActualQty
				if err := tx.Save(&product).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses pengajuan SO: " + err.Error()})
		return
	}
	
	msg := "Berkas berkuitansi Stock Opname berhasil diteruskan ke laptop Owner !"
	if status == "APPROVED" {
		msg = "Audit Stock Opname berhasil dieksekusi, stok master cabang diperbarui!"
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *RetailHandler) ApproveStockOpname(c *gin.Context) {
	opnameID := c.Param("id")
	userRole := strings.ToLower(c.MustGet("role").(string))
	if userRole != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hak otorisasi ditolak! Hanya akun Owner yang memiliki akses persetujuan."})
		return
	}

	file, errFile := c.FormFile("bukti_bar")
	if errFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dokumen fisik PDF Berita Acara (BAR TTD) wajib diunggah !"})
		return
	}

	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")
	fileSrc, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membedah berkas file PDF"})
		return
	}
	defer fileSrc.Close()

	remotePath := fmt.Sprintf("audit/so_%s_bar", opnameID)
	urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, "application/pdf", bucketName, remotePath)
	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Koneksi cloud storage storage terputus, gagal simpan dokumen BAR"})
		return
	}

	db := h.Repo.GetDB()
err = db.Transaction(func(tx *gorm.DB) error {
    var so domain.StockOpname
    
    // 🚀 FIX: Eksplisit sebut nama kolomnya biar GORM ga salah deteksi tipe data ULID/String bray!
    // Sesuaikan kolomnya ya, pakai "public_id = ?" atau "id = ?" tergantung skema tabel SO lu.
    if err := tx.Preload("Details").Clauses(clause.Locking{Strength: "UPDATE"}).First(&so, "public_id = ?", opnameID).Error; err != nil {
        return fmt.Errorf("berkas ID SO tidak ditemukan atau gagal dikunci")
    }
    
    if so.Status == "APPROVED" {
        return nil
    }

    if err := tx.Model(&so).Updates(map[string]interface{}{"status": "APPROVED", "bukti_bar": urlResult}).Error; err != nil {
        return err
    }

    for _, detail := range so.Details {
        var product models.Product
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, "id = ? AND store_id = ?", detail.ProductID, so.StoreID).Error; err != nil {
            return err
        }
        product.Stok = detail.ActualQty
        if err := tx.Save(&product).Error; err != nil {
            return err
        }
    }
    return nil
})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal eksekusi persetujuan SO : " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Audit berhasil disetujui, stok master diperbarui di Cloud!"})
}

// ==========================================
// 🚀 HISTORI & VALIDASI KELAYAKAN
// ==========================================

type AuditCompareResult struct {
	SO    domain.StockOpname       `json:"so"`
	Klaim *domain.StockAdjustment `json:"klaim"`
}

func (h *RetailHandler) GetStockOpnameHistory(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	db := h.Repo.GetDB()
	var soHistory []domain.StockOpname
	err := db.Where("store_id = ?", storeID).Preload("Details.Product").Order("created_at desc").Find(&soHistory).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyisir riwayat dokumen berkas opname"})
		return
	}

	var klaimHistory []domain.StockAdjustment
	db.Where("store_id = ?", storeID).Preload("Details.Product").Find(&klaimHistory)

	var results []AuditCompareResult
	for _, so := range soHistory {
		compare := AuditCompareResult{SO: so}
		for _, klaim := range klaimHistory {
			if klaim.CreatedAt.Month() == so.CreatedAt.Month() && klaim.CreatedAt.Year() == so.CreatedAt.Year() {
				klaimCopy := klaim
				compare.Klaim = &klaimCopy
				break
			}
		}
		results = append(results, compare)
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}

func (h *RetailHandler) GetLastSOStatus(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	var lastSO domain.StockOpname
	db := h.Repo.GetDB()
	result := db.Where("store_id = ? AND status = ?", storeID, "APPROVED").Order("created_at desc").First(&lastSO)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"last_so_date": nil, "has_claimed": false})
		return
	}

	var claimCount int64
	db.Model(&domain.StockAdjustment{}).Where("store_id = ? AND created_at >= ?", storeID, lastSO.CreatedAt).Count(&claimCount)
	c.JSON(http.StatusOK, gin.H{"last_so_date": lastSO.CreatedAt, "has_claimed": claimCount > 0})
}

func (h *RetailHandler) GetLastSOMinusItems(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	db := h.Repo.GetDB()
	var lastSO domain.StockOpname
	if err := db.Where("store_id = ? AND status = ?", storeID, "APPROVED").Order("created_at desc").First(&lastSO).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": []domain.StockOpnameDetail{}})
		return
	}

	var minusDetails []domain.StockOpnameDetail
	db.Where("opname_id = ? AND selisih < 0", lastSO.ID).Preload("Product").Find(&minusDetails)
	c.JSON(http.StatusOK, gin.H{"data": minusDetails})
}

// ==========================================
// 🚀 KLAIM BARANG ADJUSTMENT SYSTEM UNIT
// ==========================================

func (h *RetailHandler) SubmitKlaimBarang(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }
	switch v := userIDRaw.(type) { case float64: userID = uint(v); case uint: userID = v; case int: userID = uint(v) }

	var req KlaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format berkas klaim tidak sesuai !"})
		return
	}

	db := h.Repo.GetDB()
	var lastSO domain.StockOpname
	if err := db.Where("store_id = ? AND status = ?", storeID, "APPROVED").Order("created_at desc").First(&lastSO).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dilarang klaim! Toko Anda belum pernah melakukan audit Stock Opname bulanan."})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		adj := domain.StockAdjustment{
			PublicID:  utils.GenerateULID(), 
			StoreID:   storeID,
			UserID:    userID,
			Notes:     req.Notes,
			Status:    "PENDING_APPROVAL",
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&adj).Error; err != nil {
			return err
		}

		for _, item := range req.Items {
			var soDetail domain.StockOpnameDetail
			err := tx.Where("opname_id = ? AND product_id = ? AND selisih < 0", lastSO.ID, item.ProductID).First(&soDetail).Error
			if err != nil {
				return fmt.Errorf("produk ID %d ditolak! Barang ini tidak terdaftar sebagai produk hilang minus di SO terakhir", item.ProductID)
			}

			// 🛡️ SECURITY FIX: Hitung kumulatif qty yang SUDAH PERNAH diajukan sebelumnya (Pending/Approved)
			var totalClaimedBefore int64
			err = tx.Model(&domain.StockAdjustmentDetail{}).
				Joins("JOIN stock_adjustments ON stock_adjustments.id = stock_adjustment_details.adjustment_id").
				Where("stock_adjustments.store_id = ? AND stock_adjustments.created_at >= ? AND stock_adjustment_details.product_id = ?", 
					storeID, lastSO.CreatedAt, item.ProductID).
				Select("COALESCE(SUM(stock_adjustment_details.qty), 0)").
				Row().Scan(&totalClaimedBefore)
			
			if err != nil {
				return fmt.Errorf("gagal memvalidasi histori limit klaim produk ID %d", item.ProductID)
			}

			allowedMaxLimit := int(math.Abs(float64(soDetail.Selisih)))
			sisaKuotaKlaim := allowedMaxLimit - int(totalClaimedBefore)

			if item.Qty > sisaKuotaKlaim {
				return fmt.Errorf("klaim ditolak! Produk ID %d melebihi batas sisa kuota kehilangan (Sisa Kuota: %d PCS, Sudah Diklaim Sebelumnya: %d PCS)", 
					item.ProductID, sisaKuotaKlaim, totalClaimedBefore)
			}

			adjItem := domain.StockAdjustmentDetail{
				AdjustmentID: adj.ID, 
				ProductID:    item.ProductID,
				Qty:          item.Qty,
				Alasan:       item.Alasan,
			}
			if err := tx.Create(&adjItem).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Klaim penemuan penyesuaian dana barang berhasil diteruskan ke laptop Owner !"})
}

func (h *RetailHandler) GetStockAdjustmentHistory(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	var history []domain.StockAdjustment
	db := h.Repo.GetDB()
	err := db.Where("store_id = ?", storeID).Preload("Details.Product").Order("created_at desc").Find(&history).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengumpulkan riwayat kuitansi klaim"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": history})
}

func (h *RetailHandler) ApproveStockAdjustment(c *gin.Context) {
	adjustmentID := c.Param("id")
	userRole := strings.ToLower(c.MustGet("role").(string))
	if userRole != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses otorisasi ditolak! Menu khusus akun Owner."})
		return
	}

	file, errFile := c.FormFile("bukti_bar")
	if errFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Berkas kuitansi cetak PDF Berita Acara Klaim wajib diupload!"})
		return
	}

	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")
	fileSrc, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuka data lampiran PDF"})
		return
	}
	defer fileSrc.Close()

	remotePath := fmt.Sprintf("audit/claim_%s_bar", adjustmentID)
	urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, "application/pdf", bucketName, remotePath)
	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengamankan berkas BAR Klaim ke cloud storage"})
		return
	}

	db := h.Repo.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		var adjustment domain.StockAdjustment
		if err := tx.Preload("Details").Clauses(clause.Locking{Strength: "UPDATE"}).First(&adjustment, adjustmentID).Error; err != nil {
			return fmt.Errorf("berkas kode ID Klaim tidak valid")
		}
		if adjustment.Status == "APPROVED" {
			return nil
		}

		// 🛠️ ALIGNMENT FIX: Pakai .Updates() biar hook record update_at jalan semestinya bray
		if err := tx.Model(&adjustment).Updates(map[string]interface{}{"status": "APPROVED", "bukti_bar": urlResult}).Error; err != nil {
			return err
		}

		for _, detail := range adjustment.Details {
			var product models.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, "id = ? AND store_id = ?", detail.ProductID, adjustment.StoreID).Error; err != nil {
				return err
			}
			product.Stok += detail.Qty
			if err := tx.Save(&product).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyetujui kuitansi klaim : " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Klaim berhasil disetujui, jumlah stok fisik rak bertambah di Cloud!"})
}