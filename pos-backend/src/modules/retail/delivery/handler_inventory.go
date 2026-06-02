package delivery

import (
	"fmt"
	"net/http"
	"time"
	"strconv"

	"pos-backend/src/modules/retail/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 🚀 STOCK OPNAME HANDLERS
type SOItemInput struct {
	ProductID uint   `json:"product_id"`
	SystemQty int    `json:"system_qty"`
	ActualQty int    `json:"actual_qty"`
	Selisih   int    `json:"selisih"`
	Alasan    string `json:"alasan"`
}

type StockOpnameInput struct {
	Notes  string        `json:"notes"`
	Status string        `json:"status"` // 'APPROVED' atau 'PENDING_APPROVAL'
	Items  []SOItemInput `json:"items"`
}

type KlaimItemInput struct {
	ProductID uint   `json:"product_id"`
	Qty       int    `json:"qty"`
	Alasan    string `json:"alasan"` // Alasan ketemunya di mana
}

type KlaimRequest struct {
	Notes string           `json:"notes"` // Misal: "Klaim Barang Nyempil"
	Items []KlaimItemInput `json:"items"`
}

func (h *RetailHandler) CreateStockOpname(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	userID := uint(c.MustGet("user_id").(float64))
	userRole := c.MustGet("role").(string)

	var input StockOpnameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid!"})
		return
	}

	// 🚀 LOGIKA GEMBOK 1 BULAN 1 KALI
	now := time.Now()
	isAlreadyDone := h.Repo.CheckStockOpnameThisMonth(storeID, int(now.Month()), now.Year())
	if isAlreadyDone {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sistem Terkunci! Stock Opname / Audit hanya bisa dilakukan 1x dalam bulan yang sama."})
		return
	}

	status := "PENDING_APPROVAL"
	if userRole == "owner" {
		status = "APPROVED"
	}

	db := h.Repo.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		so := domain.StockOpname{
			StoreID:   storeID,
			UserID:    userID,
			Notes:     input.Notes,
			Status:    status,
			CreatedAt: time.Now(),
		}

		if err := h.Repo.CreateStockOpname(tx, &so); err != nil {
			return err
		}

		for _, item := range input.Items {
			detail := domain.StockOpnameDetail{
				OpnameID:  so.ID,
				ProductID: item.ProductID,
				SystemQty: item.SystemQty,
				ActualQty: item.ActualQty,
				Selisih:   item.Selisih,
				Alasan:    item.Alasan,
			}

			if err := h.Repo.CreateStockOpnameDetail(tx, &detail); err != nil {
				return err
			}

			if status == "APPROVED" {
				product, err := h.Repo.GetProductByID(tx, item.ProductID, storeID)
				if err != nil {
					return err
				}

				product.Stok = item.ActualQty
				if err := h.Repo.SaveProduct(tx, product); err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal proses SO: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Audit berhasil diajukan"})
}

func (h *RetailHandler) ApproveStockOpname(c *gin.Context) {
	opnameID := c.Param("id") // Ambil ID dari URL /api/retail/stock-opname/:id/approve

	// Pastikan cuma Owner yang bisa akses
	userRole := c.MustGet("role").(string)
	if userRole != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Owner yang bisa menyetujui audit!"})
		return
	}

	// 🚀 1. TANGKAP FILE PDF DARI VUE DULU
	file, errFile := c.FormFile("bukti_bar")
	var filePath string

	if errFile == nil {
		// Bikin nama unik
		filename := fmt.Sprintf("BAR_SO_%s_%d.pdf", opnameID, time.Now().Unix())
		filePath = "uploads/bar/" + filename

		// Simpan file fisik ke folder server (Pastikan folder 'uploads/bar' udah ada di root project lu)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file PDF Berita Acara"})
			return
		}
	} else {
		// Tolak kalau owner nakal mau approve tanpa PDF
		c.JSON(http.StatusBadRequest, gin.H{"error": "File PDF Berita Acara (TTD) wajib diupload!"})
		return
	}

	db := h.Repo.GetDB()

	// 🚀 2. JALANKAN TRANSAKSI DATABASE LU YANG UDAH CAKEP INI
	err := db.Transaction(func(tx *gorm.DB) error {
		var so domain.StockOpname
		if err := tx.Preload("Details").First(&so, opnameID).Error; err != nil {
			return err
		}

		if so.Status == "APPROVED" {
			return nil // Sudah pernah di-approve
		}

		// Update status SO jadi APPROVED sekaligus simpan link PDF-nya
		if err := tx.Model(&so).Updates(map[string]interface{}{
			"status":    "APPROVED",
			"bukti_bar": filePath, // 🚀 SIMPAN JEJAK PDF KE DATABASE
		}).Error; err != nil {
			return err
		}

		// Loop item dan update stok Master Product
		for _, detail := range so.Details {
			product, err := h.Repo.GetProductByID(tx, detail.ProductID, so.StoreID)
			if err != nil {
				return err
			}

			// Update stok master dengan hasil fisik (ActualQty) dari SO
			product.Stok = detail.ActualQty
			if err := h.Repo.SaveProduct(tx, product); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal approve SO: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Audit berhasil disetujui, stok master diupdate, & BAR tersimpan!"})
}

// 🚀 BIKIN STRUCT BARU BUAT NAMPUNG DATA JODOHAN (COMPARE)
type AuditCompareResult struct {
	SO    domain.StockOpname      `json:"so"`
	Klaim *domain.StockAdjustment `json:"klaim"` // Pakai pointer (*) biar bisa bernilai null kalau belum ada klaim
}

func (h *RetailHandler) GetStockOpnameHistory(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	db := h.Repo.GetDB()

	// 1. Tarik semua riwayat SO Asli
	var soHistory []domain.StockOpname
	err := db.Where("store_id = ?", storeID).
		Preload("Details.Product").
		Order("created_at desc").
		Find(&soHistory).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data riwayat opname"})
		return
	}

	// 2. Tarik semua riwayat Klaim Barang (Kalau ada)
	var klaimHistory []domain.StockAdjustment
	db.Where("store_id = ?", storeID).
		Preload("Details.Product").
		Find(&klaimHistory)

	// 3. 🚀 PROSES KAWIN SILANG (COMPARE) BERDASARKAN BULAN & TAHUN
	var results []AuditCompareResult
	for _, so := range soHistory {
		compare := AuditCompareResult{
			SO: so,
		}

		// Cari klaim yang terjadi di bulan dan tahun yang sama dengan SO ini
		for _, klaim := range klaimHistory {
			if klaim.CreatedAt.Month() == so.CreatedAt.Month() && klaim.CreatedAt.Year() == so.CreatedAt.Year() {
				klaimCopy := klaim
				compare.Klaim = &klaimCopy
				break // Langsung stop pencarian kalau udah ketemu jodohnya
			}
		}

		results = append(results, compare)
	}

	// Balikin data yang udah digabung ke Vue!
	c.JSON(http.StatusOK, gin.H{"data": results})
}

func (h *RetailHandler) GetLastSOStatus(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64)) // 🚀 Filter by Toko

	var lastSO domain.StockOpname
	db := h.Repo.GetDB()

	// Cari SO terakhir di toko ini yang statusnya APPROVED
	result := db.Where("store_id = ? AND status = ?", storeID, "APPROVED").
		Order("created_at desc").
		First(&lastSO)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"last_so_date": nil, "has_claimed": false})
		return
	}

	// 🚀 CEK APAKAH SUDAH ADA KLAIM SETELAH SO INI
	var claimCount int64
	db.Model(&domain.StockAdjustment{}).
		Where("store_id = ? AND created_at >= ?", storeID, lastSO.CreatedAt).
		Count(&claimCount)

	// Balikin tanggal SO dan status jatah klaimnya
	c.JSON(http.StatusOK, gin.H{
		"last_so_date": lastSO.CreatedAt,
		"has_claimed":  claimCount > 0, // True kalau udah pernah klaim
	})
}

// 🚀 ENDPOINT SAKTI: Narik barang yang minus di SO Terakhir
func (h *RetailHandler) GetLastSOMinusItems(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	db := h.Repo.GetDB()

	var lastSO domain.StockOpname
	if err := db.Where("store_id = ? AND status = ?", storeID, "APPROVED").Order("created_at desc").First(&lastSO).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": []domain.StockOpnameDetail{}})
		return
	}

	var minusDetails []domain.StockOpnameDetail
	// Tarik detail yang selisihnya kurang dari 0 (Barang Hilang)
	db.Where("opname_id = ? AND selisih < 0", lastSO.ID).Preload("Product").Find(&minusDetails)

	c.JSON(http.StatusOK, gin.H{"data": minusDetails})
}

func (h *RetailHandler) SubmitKlaimBarang(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	userID := uint(c.MustGet("user_id").(float64))

	var req KlaimRequest
	// 🚀 Karena gak ada foto, tetep pake ShouldBindJSON (Sangat simpel!)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	db := h.Repo.GetDB()
	tx := db.Begin()

	// Simpan ke tabel adjustment dengan status PENDING_APPROVAL
	adj := domain.StockAdjustment{
		StoreID:   storeID,
		UserID:    userID,
		Notes:     req.Notes,
		Status:    "PENDING_APPROVAL",
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&adj).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan klaim"})
		return
	}

	for _, item := range req.Items {
		adjItem := domain.StockAdjustmentDetail{
			AdjustmentID: adj.ID,
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
	storeID := uint(c.MustGet("store_id").(float64))

	var history []domain.StockAdjustment
	db := h.Repo.GetDB()

	// Tarik data klaim, preload detail barang dan data produknya
	err := db.Where("store_id = ?", storeID).
		Preload("Details.Product").
		Order("created_at desc").
		Find(&history).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat klaim"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}

// 2. Fungsi Owner Ngetok Palu Approve Klaim Barang
func (h *RetailHandler) ApproveStockAdjustment(c *gin.Context) {
	adjustmentID := c.Param("id") // Ambil ID dari URL

	// Pastikan cuma Owner yang bisa akses
	userRole := c.MustGet("role").(string)
	if userRole != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Owner yang bisa menyetujui klaim!"})
		return
	}

	// 🚀 1. TANGKAP FILE PDF DARI VUE (KHUSUS KLAIM)
	file, errFile := c.FormFile("bukti_bar")
	var filePath string

	if errFile == nil {
		// Bikin nama unik (Pake BAR_KLAIM biar beda sama SO awal)
		filename := fmt.Sprintf("BAR_KLAIM_%s_%d.pdf", adjustmentID, time.Now().Unix())
		filePath = "uploads/bar/" + filename

		// Simpan file fisik ke folder server
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file PDF Berita Acara Klaim"})
			return
		}
	} else {
		// Tolak kalau owner nakal mau approve tanpa PDF
		c.JSON(http.StatusBadRequest, gin.H{"error": "File PDF Berita Acara Klaim wajib diupload!"})
		return
	}

	db := h.Repo.GetDB()

	// 🚀 2. JALANKAN TRANSAKSI DATABASE
	err := db.Transaction(func(tx *gorm.DB) error {
		var adjustment domain.StockAdjustment
		if err := tx.Preload("Details").First(&adjustment, adjustmentID).Error; err != nil {
			return err
		}

		if adjustment.Status == "APPROVED" {
			return nil // Sudah pernah di-approve
		}

		// Update status Klaim jadi APPROVED sekaligus simpan link PDF-nya
		if err := tx.Model(&adjustment).UpdateColumns(map[string]interface{}{
			"status":    "APPROVED",
			"bukti_bar": filePath, // 🚀 SIMPAN JEJAK PDF KLAIM KE DATABASE
		}).Error; err != nil {
			return err
		}

		// Loop item dan TAMBAH stok Master Product (Karena ini barang ketemu)
		for _, detail := range adjustment.Details {
			product, err := h.Repo.GetProductByID(tx, detail.ProductID, adjustment.StoreID)
			if err != nil {
				return err
			}

			// 🚀 KLAIM = BARANG KETEMU = STOK DITAMBAH
			product.Stok += detail.Qty
			if err := h.Repo.SaveProduct(tx, product); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal approve Klaim: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Klaim berhasil disetujui, stok bertambah, & BAR Klaim tersimpan!"})
}


// ===============================
// 🚀 RETUR BARANG HANDLERS
// ===============================
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
	storeID := uint(c.MustGet("store_id").(float64))
	userID := uint(c.MustGet("user_id").(float64))

	var input ReturnInputBatch
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data keranjang retur tidak valid!"})
		return
	}

	returnNo := fmt.Sprintf("RET-%s-%d", time.Now().Format("060102150405"), userID)

	db := h.Repo.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var newReturns []domain.ProductReturn

	for _, item := range input.Items {
		product, err := h.Repo.GetProductByID(tx, item.ProductID, storeID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Ada produk yang tidak ditemukan di toko ini!"})
			return
		}

		if product.Stok < item.Qty {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stok %s tidak mencukupi (Sisa: %d)!", product.NamaProduk, product.Stok)})
			return
		}

		if err := h.Repo.UpdateProductStokExpr(tx, item.ProductID, item.Qty); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memotong stok produk!"})
			return
		}

		newReturns = append(newReturns, domain.ProductReturn{
			ReturnNo:  returnNo,
			StoreID:   storeID,
			ProductID: item.ProductID,
			UserID:    userID,
			Qty:       item.Qty,
			Alasan:    item.Alasan,
			Catatan:   item.Catatan,
		})
	}

	if err := h.Repo.CreateProductReturns(tx, newReturns); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat log retur batch!"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Berita Acara Retur berhasil diproses!", "return_no": returnNo})
}

func (h *RetailHandler) GetReturns(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	limit, offset := 0, 0
	if pageStr != "" && limitStr != "" {
		p, _ := strconv.Atoi(pageStr)
		l, _ := strconv.Atoi(limitStr)
		limit = l
		offset = (p - 1) * limit
	}

	returns, totalItems, err := h.Repo.GetReturnsHistory(storeID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data retur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Data retur berhasil dimuat!",
		"total_items": totalItems,
		"data":        returns,
	})
}


// ==============================
// 🚀 PURCHASE / LPB HANDLERS
// ==============================
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
	storeID := uint(c.MustGet("store_id").(float64))
	userID := uint(c.MustGet("user_id").(float64))

	var input PurchaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data LPB tidak valid!"})
		return
	}

	if len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keranjang penerimaan kosong!"})
		return
	}

	// 1. Susun Struktur Header Transaksi Pembelian
	purchase := domain.Purchase{
		StoreID:      storeID,
		UserID:       userID,
		SupplierName: input.SupplierName,
		NoFaktur:     input.NoFaktur,
		TotalItem:    len(input.Items),
	}

	// 2. Susun Detail Item Pembelian
	var details []domain.PurchaseDetail
	for _, item := range input.Items {
		details = append(details, domain.PurchaseDetail{
			ProductID:  item.ProductID,
			QtyMasuk:   item.QtyMasuk,
			HargaModal: item.HargaModal,
		})
	}
	purchase.Details = details

	db := h.Repo.GetDB()

	// 3. 🚀 EKSEKUSI PENYIMPANAN + MOVING AVERAGE COST KE REPOSITORY
	// Semua kalkulasi stok akhir dan perubahan harga master diselesaikan secara Atomik oleh GORM di sini.
	if err := h.Repo.CreatePurchaseWithMovingAverage(db, &purchase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses LPB dan HPP: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Penerimaan Barang berhasil! Stok & Modal HPP telah di-update."})
}