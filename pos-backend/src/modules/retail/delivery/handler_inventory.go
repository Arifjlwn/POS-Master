package delivery

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"pos-backend/src/modules/retail/domain"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ==========================================
// 🚀 STOCK OPNAME (AUDIT STOK) HANDLERS
// ==========================================

type SOItemInput struct {
	ProductID uint   `json:"product_id"`
	SystemQty int    `json:"system_qty"`
	ActualQty int    `json:"actual_qty"`
	Selisih   int    `json:"selisih"`
	Alasan    string `json:"alasan"`
}

type StockOpnameInput struct {
	Notes  string        `json:"notes"`
	Status string        `json:"status"`
	Items  []SOItemInput `json:"items"`
}

type KlaimItemInput struct {
	ProductID uint   `json:"product_id"`
	Qty       int    `json:"qty"`
	Alasan    string `json:"alasan"`
}

type KlaimRequest struct {
	Notes string           `json:"notes"`
	Items []KlaimItemInput `json:"items"`
}

func (h *RetailHandler) CreateStockOpname(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id"); userIDRaw, _ := c.Get("user_id"); userRole := c.MustGet("role").(string)
	var storeID, userID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }
	switch v := userIDRaw.(type) { case float64: userID = uint(v); case uint: userID = v; case int: userID = uint(v) }

	var input StockOpnameInput
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid!"}); return }

	now := time.Now()
	if h.Repo.CheckStockOpnameThisMonth(storeID, int(now.Month()), now.Year()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sistem Terkunci! Stock Opname hanya bisa dilakukan 1x dalam bulan yang sama."})
		return
	}

	status := "PENDING_APPROVAL"
	if userRole == "owner" { status = "APPROVED" }

	db := h.Repo.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		so := domain.StockOpname{
			PublicID:  utils.GenerateULID(), // 🚀 HYBRID MASTER: Injeksi ULID di level Header SO
			StoreID:   storeID,
			UserID:    userID,
			Notes:     input.Notes,
			Status:    status,
			CreatedAt: now,
		}
		if err := h.Repo.CreateStockOpname(tx, &so); err != nil { return err }

		for _, item := range input.Items {
			product, err := h.Repo.GetProductByID(tx, item.ProductID, storeID)
			if err != nil { return err }

			// Kalkulasi Nilai Kerugian/Keuntungan Keuangan Uang Berdasarkan Selisih Fisik Barang
			nilaiKerugian := float64(item.Selisih) * product.HargaModal

			detail := domain.StockOpnameDetail{
				OpnameID:  so.ID, // 🚀 HYBRID OPTIMAL: Detail bersih tanpa PublicID ULID acak
				ProductID: item.ProductID,
				SystemQty: item.SystemQty,
				ActualQty: item.ActualQty,
				Selisih:   item.Selisih,
				NilaiUang: nilaiKerugian,
				Alasan:    item.Alasan,
			}
			if err := h.Repo.CreateStockOpnameDetail(tx, &detail); err != nil { return err }

			if status == "APPROVED" {
				product.Stok = item.ActualQty
				if err := h.Repo.SaveProduct(tx, product); err != nil { return err }
			}
		}
		return nil
	})

	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal proses SO: " + err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Audit stok gudang berhasil diajukan!"})
}

func (h *RetailHandler) ApproveStockOpname(c *gin.Context) {
	opnameID := c.Param("id"); userRole := c.MustGet("role").(string)
	if userRole != "owner" { c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Owner yang bisa menyetujui audit!"}); return }

	file, errFile := c.FormFile("bukti_bar")
	if errFile != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "File PDF Berita Acara (BAR TTD) wajib diupload!"}); return }

	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")
	fileSrc, err := file.Open(); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file PDF"}); return }
	defer fileSrc.Close()

	// 🚀 ENTERPRISE CLOUD STORAGE: Lempar Berita Acara PDF langsung ke Supabase Cloud Storage bray
	remotePath := fmt.Sprintf("audit/so_%s_bar", opnameID)
	urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, "application/pdf", bucketName, remotePath)
	if errUpload != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload BAR ke cloud storage"}); return }

	db := h.Repo.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		var so domain.StockOpname
		if err := tx.Preload("Details").First(&so, opnameID).Error; err != nil { return err }
		if so.Status == "APPROVED" { return nil }

		if err := tx.Model(&so).Updates(map[string]interface{}{"status": "APPROVED", "bukti_bar": urlResult}).Error; err != nil { return err }

		for _, detail := range so.Details {
			product, err := h.Repo.GetProductByID(tx, detail.ProductID, so.StoreID)
			if err != nil { return err }
			product.Stok = detail.ActualQty
			if err := h.Repo.SaveProduct(tx, product); err != nil { return err }
		}
		return nil
	})

	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal approve SO: " + err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Audit berhasil disetujui, stok master diperbarui di Cloud!"})
}

type AuditCompareResult struct {
	SO    domain.StockOpname      `json:"so"`
	Klaim *domain.StockAdjustment `json:"klaim"`
}

func (h *RetailHandler) GetStockOpnameHistory(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	db := h.Repo.GetDB()
	var soHistory []domain.StockOpname
	err := db.Where("store_id = ?", storeID).Preload("Details.Product").Order("created_at desc").Find(&soHistory).Error
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data riwayat opname"}); return }

	var klaimHistory []domain.StockAdjustment
	db.Where("store_id = ?", storeID).Preload("Details.Product").Find(&klaimHistory)

	var results []AuditCompareResult
	for _, so := range soHistory {
		compare := AuditCompareResult{SO: so}
		for _, klaim := range klaimHistory {
			if klaim.CreatedAt.Month() == so.CreatedAt.Month() && klaim.CreatedAt.Year() == so.CreatedAt.Year() {
				klaimCopy := klaim; compare.Klaim = &klaimCopy; break
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

	var lastSO domain.StockOpname; db := h.Repo.GetDB()
	result := db.Where("store_id = ? AND status = ?", storeID, "APPROVED").Order("created_at desc").First(&lastSO)
	if result.Error != nil { c.JSON(http.StatusOK, gin.H{"last_so_date": nil, "has_claimed": false}); return }

	var claimCount int64
	db.Model(&domain.StockAdjustment{}).Where("store_id = ? AND created_at >= ?", storeID, lastSO.CreatedAt).Count(&claimCount)
	c.JSON(http.StatusOK, gin.H{"last_so_date": lastSO.CreatedAt, "has_claimed": claimCount > 0})
}

func (h *RetailHandler) GetLastSOMinusItems(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	db := h.Repo.GetDB(); var lastSO domain.StockOpname
	if err := db.Where("store_id = ? AND status = ?", storeID, "APPROVED").Order("created_at desc").First(&lastSO).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": []domain.StockOpnameDetail{}})
		return
	}

	var minusDetails []domain.StockOpnameDetail
	db.Where("opname_id = ? AND selisih < 0", lastSO.ID).Preload("Product").Find(&minusDetails)
	c.JSON(http.StatusOK, gin.H{"data": minusDetails})
}

func (h *RetailHandler) SubmitKlaimBarang(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id"); userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }
	switch v := userIDRaw.(type) { case float64: userID = uint(v); case uint: userID = v; case int: userID = uint(v) }

	var req KlaimRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"}); return }

	db := h.Repo.GetDB(); tx := db.Begin()
	adj := domain.StockAdjustment{
		PublicID:  utils.GenerateULID(), // 🚀 HYBRID MASTER: Injeksi ULID di level Header Adjustment Klaim
		StoreID:   storeID,
		UserID:    userID,
		Notes:     req.Notes,
		Status:    "PENDING_APPROVAL",
		CreatedAt: time.Now(),
	}
	if err := tx.Create(&adj).Error; err != nil { tx.Rollback(); c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan klaim"}); return }

	for _, item := range req.Items {
		adjItem := domain.StockAdjustmentDetail{
			AdjustmentID: adj.ID, // 🚀 HYBRID OPTIMAL: Detail bersih tanpa ULID string sampah
			ProductID:    item.ProductID,
			Qty:          item.Qty,
			Alasan:       item.Alasan,
		}
		tx.Create(&adjItem)
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Klaim barang temuan berhasil diajukan ke Owner!"})
}

func (h *RetailHandler) GetStockAdjustmentHistory(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	var history []domain.StockAdjustment; db := h.Repo.GetDB()
	err := db.Where("store_id = ?", storeID).Preload("Details.Product").Order("created_at desc").Find(&history).Error
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat klaim"}); return }
	c.JSON(http.StatusOK, gin.H{"data": history})
}

func (h *RetailHandler) ApproveStockAdjustment(c *gin.Context) {
	adjustmentID := c.Param("id"); userRole := c.MustGet("role").(string)
	if userRole != "owner" { c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Owner yang bisa menyetujui klaim!"}); return }

	file, errFile := c.FormFile("bukti_bar")
	if errFile != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "File PDF Berita Acara Klaim wajib diupload!"}); return }

	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")
	fileSrc, err := file.Open(); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file PDF"}); return }
	defer fileSrc.Close()

	// 🚀 ENTERPRISE CLOUD STORAGE: Upload Berita Acara Klaim PDF ke Supabase
	remotePath := fmt.Sprintf("audit/claim_%s_bar", adjustmentID)
	urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, "application/pdf", bucketName, remotePath)
	if errUpload != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload BAR Klaim ke cloud storage"}); return }

	db := h.Repo.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		var adjustment domain.StockAdjustment
		if err := tx.Preload("Details").First(&adjustment, adjustmentID).Error; err != nil { return err }
		if adjustment.Status == "APPROVED" { return nil }

		if err := tx.Model(&adjustment).UpdateColumns(map[string]interface{}{"status": "APPROVED", "bukti_bar": urlResult}).Error; err != nil { return err }

		for _, detail := range adjustment.Details {
			product, err := h.Repo.GetProductByID(tx, detail.ProductID, adjustment.StoreID)
			if err != nil { return err }
			product.Stok += detail.Qty
			if err := h.Repo.SaveProduct(tx, product); err != nil { return err }
		}
		return nil
	})

	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal approve Klaim: " + err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Klaim berhasil disetujui, stok bertambah di Cloud!"})
}

// ==========================================
// 🚀 RETUR BARANG (PRODUCT RETURN) HANDLERS
// ==========================================

type ReturnItem struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Qty       int    `json:"qty" binding:"required,gt=0"`
	Alasan    string `json:"alasan" binding:"required"`
	Catatan   string `json:"catatan"`
}

type ReturnInputBatch struct {
	Items []ReturnItem `json:"items" binding:"required,min=1"`
}

func (h *RetailHandler) CreateReturn(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id"); userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }
	switch v := userIDRaw.(type) { case float64: userID = uint(v); case uint: userID = v; case int: userID = uint(v) }

	var input ReturnInputBatch
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Data keranjang retur tidak valid!"}); return }

	returnNo := fmt.Sprintf("RET-%s-%d", time.Now().Format("060102150405"), userID)
	db := h.Repo.GetDB(); tx := db.Begin()
	defer func() { if r := recover(); r != nil { tx.Rollback() } }()

	var newReturns []domain.ProductReturn
	for _, item := range input.Items {
		product, err := h.Repo.GetProductByID(tx, item.ProductID, storeID)
		if err != nil { tx.Rollback(); c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan!"}); return }
		if product.Stok < item.Qty { tx.Rollback(); c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stok %s tidak cukup (Sisa: %d)!", product.NamaProduk, product.Stok)}); return }

		if err := h.Repo.UpdateProductStokExpr(tx, item.ProductID, item.Qty); err != nil { tx.Rollback(); c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal potong stok!"}); return }

		newReturns = append(newReturns, domain.ProductReturn{
			PublicID:  utils.GenerateULID(), // 🚀 HYBRID FLAT: Wajib di-inject ULID tiap baris karena model flat table
			ReturnNo:  returnNo,
			StoreID:   storeID,
			ProductID: item.ProductID,
			UserID:    userID,
			Qty:       item.Qty,
			Alasan:    item.Alasan,
			Catatan:   item.Catatan,
		})
	}

	if err := h.Repo.CreateProductReturns(tx, newReturns); err != nil { tx.Rollback(); c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat retur!"}); return }
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Berita Acara Retur berhasil diproses!", "return_no": returnNo})
}

func (h *RetailHandler) GetReturns(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	pageStr := c.Query("page"); limitStr := c.Query("limit")
	limit, offset := 0, 0
	if pageStr != "" && limitStr != "" {
		p, _ := strconv.Atoi(pageStr); l, _ := strconv.Atoi(limitStr)
		limit = l; offset = (p - 1) * limit
	}

	returns, totalItems, err := h.Repo.GetReturnsHistory(storeID, limit, offset)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data retur"}); return }
	c.JSON(http.StatusOK, gin.H{"message": "Data retur berhasil dimuat!", "total_items": totalItems, "data": returns})
}

// ==========================================
// 🚀 PURCHASE / LPB (PENERIMAAN BARANG) HANDLERS
// ==========================================

type PurchaseInput struct {
	SupplierName string `json:"supplier_name"`
	NoFaktur     string `json:"no_faktur"`
	Items        []struct {
		ProductID  uint    `json:"product_id"`
		QtyMasuk   int     `json:"qty_masuk"`
		HargaModal float64 `json:"harga_modal"`
	} `json:"items"`
}

func (h *RetailHandler) CreateLPB(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id"); userIDRaw, _ := c.Get("user_id")
	var storeID, userID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }
	switch v := userIDRaw.(type) { case float64: userID = uint(v); case uint: userID = v; case int: userID = uint(v) }

	var input PurchaseInput
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Data LPB tidak valid!"}); return }
	if len(input.Items) == 0 { c.JSON(http.StatusBadRequest, gin.H{"error": "Keranjang penerimaan kosong!"}); return }

	var totalHargaFaktur float64
	var details []domain.PurchaseDetail

	for _, item := range input.Items {
		subTotalItem := float64(item.QtyMasuk) * item.HargaModal
		totalHargaFaktur += subTotalItem

		details = append(details, domain.PurchaseDetail{
			ProductID:  item.ProductID, // 🚀 HYBRID OPTIMAL: Detail bersih dari ULID acak
			QtyMasuk:   item.QtyMasuk,
			HargaModal: item.HargaModal,
			SubTotal:   subTotalItem,
		})
	}

	purchase := domain.Purchase{
		PublicID:     utils.GenerateULID(), // 🚀 HYBRID MASTER: Injeksi ULID di level Header LPB Faktur
		StoreID:      storeID,
		UserID:       userID,
		SupplierName: input.SupplierName,
		NoFaktur:     input.NoFaktur,
		TotalItem:    len(input.Items),
		TotalHarga:   totalHargaFaktur, // Sinkronisasi kalkulasi nominal audit keuangan masuk/keluar bray
		StatusBayar:  "LUNAS",
		Details:      details,
	}

	db := h.Repo.GetDB()
	if err := h.Repo.CreatePurchaseWithMovingAverage(db, &purchase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses LPB dan HPP: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Penerimaan Barang berhasil! Stok & Modal HPP telah di-update via Moving Average."})
}