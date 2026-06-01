package delivery

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"pos-backend/models"
	src "pos-backend/src/core/config"
	"pos-backend/src/modules/retail/domain"
	"pos-backend/src/modules/retail/repository"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type RetailHandler struct {
	Repo repository.RetailRepository
}

func NewRetailHandler(repo repository.RetailRepository) *RetailHandler {
	return &RetailHandler{Repo: repo}
}

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

// 🚀 RETUR HANDLERS
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

// 🚀 PURCHASE / LPB HANDLERS
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

type AbsenInput struct {
	UserID uint   `json:"user_id" binding:"required"`
	Jenis  string `json:"jenis" binding:"required"`
	Foto   string `json:"foto" binding:"required"`
}

func (h *RetailHandler) StoreAttendance(c *gin.Context) {
	storeIDRaw, exists := c.Get("store_id")
	if !exists || storeIDRaw == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Toko tidak terdeteksi! Pastikan akun sudah terhubung."})
		return
	}

	var storeID uint
	if val, ok := storeIDRaw.(float64); ok {
		storeID = uint(val)
	} else if val, ok := storeIDRaw.(uint); ok {
		storeID = val
	}

	var input AbsenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid!"})
		return
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	today := now.Format("2006-01-02")
	nowTime := now.Format("15:04:05")

	attendance, err := h.Repo.GetAttendanceToday(input.UserID, today)

	if input.Jenis == "Masuk" {
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah melakukan Absen Masuk hari ini!"})
			return
		}

		absen := models.Attendance{
			StoreID:   storeID,
			UserID:    input.UserID,
			Tanggal:   today,
			JamMasuk:  nowTime,
			FotoMasuk: input.Foto,
			Status:    "Hadir",
		}

		if err := h.Repo.CreateAttendance(&absen); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi masuk!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Absen Masuk Berhasil! Selamat Bekerja."})

	} else if input.Jenis == "Pulang" {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda belum melakukan Absen Masuk hari ini!"})
			return
		}

		if attendance.JamPulang != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah melakukan Absen Pulang hari ini!"})
			return
		}

		attendance.JamPulang = nowTime
		attendance.FotoPulang = input.Foto

		if err := h.Repo.SaveAttendance(attendance); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absen pulang!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Absen Pulang Berhasil! Hati-hati di jalan."})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jenis absen tidak dikenali!"})
	}
}

func (h *RetailHandler) GetAttendance(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	if val, ok := storeIDRaw.(float64); ok {
		storeID = uint(val)
	} else if val, ok := storeIDRaw.(uint); ok {
		storeID = val
	}

	tanggal := c.Query("tanggal")
	bulan := c.Query("bulan")
	tahun := c.Query("tahun")

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	todayStr := now.Format("2006-01-02")

	var prefixBulan string
	if tanggal == "" && bulan != "" && tahun != "" {
		prefixBulan = fmt.Sprintf("%s-%s-%%", tahun, bulan)
	}

	riwayat, err := h.Repo.GetAttendanceReport(storeID, tanggal, prefixBulan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik data log absensi"})
		return
	}

	db := h.Repo.GetDB()
	for i := 0; i < len(riwayat); i++ {
		if riwayat[i].JamMasuk != "" && riwayat[i].JamPulang != "" {
			riwayat[i].Status = "Hadir"
		} else if riwayat[i].JamMasuk != "" && riwayat[i].JamPulang == "" {
			if riwayat[i].Tanggal < todayStr {
				riwayat[i].Status = "Lupa Absen Pulang"
				db.Model(&riwayat[i]).Update("status", "Lupa Absen Pulang")
			} else {
				riwayat[i].Status = "Hadir"
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": riwayat})
}

func (h *RetailHandler) ExportAttendance(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	bulan := c.Query("bulan")
	tahun := c.Query("tahun")

	if bulan == "" || tahun == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulan dan tahun harus diisi!"})
		return
	}

	prefixBulan := fmt.Sprintf("%s-%s-%%", tahun, bulan)
	riwayat, err := h.Repo.GetAttendanceReport(storeID, "", prefixBulan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses ekspor laporan"})
		return
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	w.Write([]string{"Tanggal", "NIK", "Nama Karyawan", "Jam Masuk", "Jam Pulang", "Status"})

	for _, logData := range riwayat {
		nik := "-"
		if logData.User.NIK != nil {
			nik = *logData.User.NIK
		}

		jamPulang := logData.JamPulang
		if jamPulang == "" {
			jamPulang = "Belum Pulang"
		}

		w.Write([]string{
			logData.Tanggal,
			nik,
			logData.User.Name,
			logData.JamMasuk,
			jamPulang,
			logData.Status,
		})
	}
	w.Flush()

	filename := fmt.Sprintf("Laporan_Absensi_%s_%s.csv", bulan, tahun)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}

// 🤝 EMPLOYEE METHODS
func (h *RetailHandler) CreateEmployee(c *gin.Context) {
	roleOwner, _ := c.Get("role")
	if roleOwner != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Owner yang bisa mendaftarkan karyawan baru!"})
		return
	}

	storeID := uint(c.MustGet("store_id").(float64))
	name := c.PostForm("name")
	password := c.PostForm("password")
	tempatLahir := c.PostForm("tempat_lahir")
	tanggalLahir := c.PostForm("tanggal_lahir")
	noHP := c.PostForm("no_hp")
	inputRole := c.PostForm("role")

	// 🚀 PERBAIKAN: Nomor HP sekarang WAJIB DIISI!
	if name == "" || password == "" || noHP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama, Nomor HP, dan Password wajib diisi!"})
		return
	}
	if inputRole == "" {
		inputRole = "kasir"
	}

	// 🚀 FORMAT NOMOR HP JADI 628xxx SEBELUM DISIMPAN
	formattedHP := utils.FormatPhoneNumber(noHP)

	currentYear := time.Now().Format("2006")
	var newNIK string

	lastEmployee, err := h.Repo.GetLastEmployeeNIK(storeID, currentYear)
	if err != nil {
		// 🚀 JIKA BELUM ADA, MULAI DARI TAHUN + 0001 (Contoh: 20260001)
		newNIK = fmt.Sprintf("%s0001", currentYear)
	} else {
		lastNIK := *lastEmployee.NIK

		// 🚀 AMBIL 4 DIGIT TERAKHIR UNTUK URUTAN (Bukan 3 digit lagi)
		if len(lastNIK) >= 4 {
			lastSequence, _ := strconv.Atoi(lastNIK[len(lastNIK)-4:])
			newNIK = fmt.Sprintf("%s%04d", currentYear, lastSequence+1)
		} else {
			newNIK = fmt.Sprintf("%s0001", currentYear)
		}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	file, err := c.FormFile("foto")
	var fotoURL string
	if err == nil {
		newFileName := fmt.Sprintf("%s_%d%s", newNIK, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)
		c.SaveUploadedFile(file, uploadPath)
		fotoURL = "/uploads/" + newFileName
	}

	employee := models.User{
		StoreID:      &storeID,
		Name:         name,
		NIK:          &newNIK,
		Password:     string(hashedPassword),
		Role:         inputRole,
		TempatLahir:  tempatLahir,
		TanggalLahir: tanggalLahir,
		NoHP:         formattedHP, // 🚀 MASUKIN NOMOR HP YANG UDAH DIBERSIHKAN
		FotoURL:      fotoURL,
	}

	if err := h.Repo.CreateEmployee(&employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Karyawan baru berhasil didaftarkan! 🤝",
		"data":    gin.H{"nama": employee.Name, "nik": newNIK, "jabatan": employee.Role},
	})
}

func (h *RetailHandler) GetEmployees(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	employees, err := h.Repo.GetAllEmployees(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": employees})
}

func (h *RetailHandler) UpdateEmployee(c *gin.Context) {
	roleOwner, _ := c.Get("role")
	if roleOwner != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya owner yang bisa edit data tim!"})
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	employee, err := h.Repo.GetEmployeeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Karyawan tidak ditemukan!"})
		return
	}

	employee.Name = c.PostForm("name")
	employee.TempatLahir = c.PostForm("tempat_lahir")
	employee.TanggalLahir = c.PostForm("tanggal_lahir")

	// 🚀 FORMAT NOMOR HP JUGA SAAT UPDATE
	if inputHP := c.PostForm("no_hp"); inputHP != "" {
		employee.NoHP = utils.FormatPhoneNumber(inputHP)
	}

	if newRole := c.PostForm("role"); newRole != "" {
		employee.Role = newRole
	}
	if password := c.PostForm("password"); password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		employee.Password = string(hashed)
	}

	nikClean := "karyawan"
	if employee.NIK != nil {
		nikClean = *employee.NIK
	}

	if file, err := c.FormFile("foto"); err == nil {
		newFileName := fmt.Sprintf("%s_%d%s", nikClean, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)
		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
			if employee.FotoURL != "" {
				os.Remove("." + employee.FotoURL)
			}
			employee.FotoURL = "/uploads/" + newFileName
		}
	}

	if bioFile, errBio := c.FormFile("biometric_file"); errBio == nil {
		newBioName := fmt.Sprintf("%s_bio_%d%s", nikClean, time.Now().Unix(), filepath.Ext(bioFile.Filename))
		uploadBioPath := filepath.Join("uploads", newBioName)
		if err := c.SaveUploadedFile(bioFile, uploadBioPath); err == nil {
			if employee.BiometricURL != "" {
				os.Remove("." + employee.BiometricURL)
			}
			employee.BiometricURL = "/uploads/" + newBioName
		}
	}

	if err := h.Repo.SaveEmployee(employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perubahan ke database"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui! 💾", "data": employee})
}

type ProductInput struct {
	SKU            *string `form:"sku"`
	NamaProduk     string  `form:"nama_produk" binding:"required"`
	Kategori       string  `form:"kategori"`
	HargaModal     float64 `form:"harga_modal"`
	HargaJual      float64 `form:"harga_jual" binding:"required"`
	Stok           int     `form:"stok"`
	SatuanDasar    string  `form:"satuan_dasar"`
	SatuanBesar    string  `form:"satuan_besar"`
	IsiPerBesar    int     `form:"isi_per_besar"`
	HargaJualBesar float64 `form:"harga_jual_besar"`

	// 🚀 TAMBAHAN: Nangkep Form 3 Lapis (Rokok/Renteng)
	IsNestedUom      bool   `form:"is_nested_uom"`
	SatuanTengah     string `form:"satuan_tengah"`
	IsiBesarKeTengah int    `form:"isi_besar_ke_tengah"`
	IsiTengahKeDasar int    `form:"isi_tengah_ke_dasar"`
}

func (h *RetailHandler) CreateProduct(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))

	var input ProductInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Periksa kembali isian form Anda: " + err.Error()})
		return
	}
	if input.SatuanDasar == "" {
		input.SatuanDasar = "PCS"
	}

	var imagePath string
	file, errFile := c.FormFile("gambar")
	if errFile == nil {
		folderPath := "uploads/products"
		os.MkdirAll(folderPath, os.ModePerm)
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		savePath := filepath.Join(folderPath, fileName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar produk"})
			return
		}
		imagePath = "/" + savePath
	}

	product := models.Product{
		StoreID:        storeID,
		SKU:            input.SKU,
		NamaProduk:     input.NamaProduk,
		Kategori:       input.Kategori,
		HargaModal:     input.HargaModal,
		HargaJual:      input.HargaJual,
		Stok:           input.Stok,
		Gambar:         imagePath,
		SatuanDasar:    input.SatuanDasar,
		SatuanBesar:    input.SatuanBesar,
		IsiPerBesar:    input.IsiPerBesar,
		HargaJualBesar: input.HargaJualBesar,

		// 🚀 TAMBAHAN: Masukin data 3 Lapis ke Database
		IsNestedUom:      input.IsNestedUom,
		SatuanTengah:     input.SatuanTengah,
		IsiBesarKeTengah: input.IsiBesarKeTengah,
		IsiTengahKeDasar: input.IsiTengahKeDasar,
	}

	if err := h.Repo.CreateProductGlobal(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk. Barcode mungkin sudah dipakai barang lain."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Produk berhasil ditambahkan! 📦", "data": product})
}

func (h *RetailHandler) GetProducts(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	search := c.Query("search")
	category := c.Query("category")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	usePagination := false
	limit, offset := 10, 0

	if pageStr != "" || limitStr != "" {
		usePagination = true
		if pageStr == "" {
			pageStr = "1"
		}
		if limitStr == "" {
			limitStr = "10"
		}
		p, _ := strconv.Atoi(pageStr)
		l, _ := strconv.Atoi(limitStr)
		limit = l
		offset = (p - 1) * limit
	}

	products, totalItems, err := h.Repo.GetProductsCatalog(storeID, search, category, limit, offset, usePagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Katalog produk berhasil dimuat!",
		"total_items": totalItems,
		"data":        products,
	})
}

func (h *RetailHandler) UpdateProduct(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	role := c.MustGet("role").(string)

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hentikan! Cuma Owner yang boleh ubah harga/data barang."})
		return
	}

	productID, _ := strconv.Atoi(c.Param("id"))
	product, err := h.Repo.GetProductByIDSimple(uint(productID), storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan atau bukan milik toko Anda!"})
		return
	}

	product.NamaProduk = c.PostForm("nama_produk")
	if sku := c.PostForm("sku"); sku != "" {
		product.SKU = &sku
	} else {
		product.SKU = nil
	}
	product.Kategori = c.PostForm("kategori")

	if hargaModal, err := strconv.ParseFloat(c.PostForm("harga_modal"), 64); err == nil {
		product.HargaModal = hargaModal
	}
	if hargaJual, err := strconv.ParseFloat(c.PostForm("harga_jual"), 64); err == nil {
		product.HargaJual = hargaJual
	}
	if stok, err := strconv.Atoi(c.PostForm("stok")); err == nil {
		product.Stok = stok
	}
	if satuanDasar := c.PostForm("satuan_dasar"); satuanDasar != "" {
		product.SatuanDasar = satuanDasar
	}
	product.SatuanBesar = c.PostForm("satuan_besar")
	if isiPerBesar, err := strconv.Atoi(c.PostForm("isi_per_besar")); err == nil {
		product.IsiPerBesar = isiPerBesar
	} else {
		product.IsiPerBesar = 0
	}

	// Nangkep harga_jual_besar
	if hargaJualBesar, err := strconv.ParseFloat(c.PostForm("harga_jual_besar"), 64); err == nil {
		product.HargaJualBesar = hargaJualBesar
	} else {
		product.HargaJualBesar = 0
	}

	// 🚀 TAMBAHAN: Nangkep data kemasan tengah pas di-Edit
	product.IsNestedUom = c.PostForm("is_nested_uom") == "true"
	product.SatuanTengah = c.PostForm("satuan_tengah")
	if ibt, err := strconv.Atoi(c.PostForm("isi_besar_ke_tengah")); err == nil {
		product.IsiBesarKeTengah = ibt
	} else {
		product.IsiBesarKeTengah = 0
	}
	if itd, err := strconv.Atoi(c.PostForm("isi_tengah_ke_dasar")); err == nil {
		product.IsiTengahKeDasar = itd
	} else {
		product.IsiTengahKeDasar = 0
	}

	file, errFile := c.FormFile("gambar")
	if errFile == nil {
		folderPath := "uploads/products"
		os.MkdirAll(folderPath, os.ModePerm)
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		savePath := filepath.Join(folderPath, fileName)

		if err := c.SaveUploadedFile(file, savePath); err == nil {
			if product.Gambar != "" {
				os.Remove("." + product.Gambar)
			}
			product.Gambar = "/" + savePath
		}
	}

	db := h.Repo.GetDB()
	if err := h.Repo.SaveProduct(db, product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update produk. Barcode mungkin bentrok."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diubah! ✏️", "data": product})
}

func (h *RetailHandler) DeleteProduct(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	role := c.MustGet("role").(string)

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Waduh, Kasir dilarang hapus barang dari sistem!"})
		return
	}

	productID, _ := strconv.Atoi(c.Param("id"))
	product, err := h.Repo.GetProductByIDSimple(uint(productID), storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan atau bukan milik toko Anda!"})
		return
	}

	if product.Gambar != "" {
		os.Remove("." + product.Gambar)
	}
	if err := h.Repo.DeleteProductGlobal(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Barang berhasil dihapus dari gudang! 🗑️"})
}

func (h *RetailHandler) GetCategories(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	categories, err := h.Repo.GetDistinctCategories(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *RetailHandler) ExportProducts(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	products, err := h.Repo.GetAllProductsForExport(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	w.Comma = '|'

	// 🚀 1. TAMBAHIN HEADER "Harga Jual Besar"
	w.Write([]string{"SKU", "Nama Produk", "Kategori", "Harga Modal", "Harga Jual", "Stok", "Satuan Terkecil", "Satuan Tengah", "Satuan Besar", "Isi Per Besar", "Harga Jual Besar"})

	for _, p := range products {
		sku := ""
		if p.SKU != nil {
			sku = *p.SKU
		}
		w.Write([]string{
			sku, p.NamaProduk, p.Kategori,
			fmt.Sprintf("%.0f", p.HargaModal), fmt.Sprintf("%.0f", p.HargaJual),
			fmt.Sprintf("%d", p.Stok), p.SatuanDasar, p.SatuanTengah, p.SatuanBesar, fmt.Sprintf("%d", p.IsiPerBesar),
			fmt.Sprintf("%.0f", p.HargaJualBesar), // 🚀 2. MASUKIN DATANYA KE BARIS CSV
		})
	}
	w.Flush()

	c.Header("Content-Disposition", "attachment; filename = katalog_produk_pos.csv")
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}

func (h *RetailHandler) ImportProducts(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'
	_, _ = reader.Read() // Skip header

	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca isi CSV"})
		return
	}

	db := h.Repo.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, row := range records {
			sku := row[0]
			nama := row[1]
			kategori := row[2]
			modal, _ := strconv.ParseFloat(row[3], 64)
			jual, _ := strconv.ParseFloat(row[4], 64)
			stok, _ := strconv.Atoi(row[5])
			dasar := row[6]
			besar := row[7]
			isi, _ := strconv.Atoi(row[8])

			// 🚀 3. BACA KOLOM BARU (Pake pengaman if len(row) > 9 biar ga error kalo upload CSV lama)
			var jualBesar float64
			if len(row) > 9 {
				jualBesar, _ = strconv.ParseFloat(row[9], 64)
			}

			if nama == "" {
				continue
			}

			var product models.Product
			res := tx.Where("sku = ? AND store_id = ?", sku, storeID).First(&product)

			if res.Error == nil {
				// UPDATE BARANG YANG UDAH ADA
				product.NamaProduk = nama
				product.Kategori = kategori
				product.HargaModal = modal
				product.HargaJual = jual
				product.Stok = stok
				product.SatuanDasar = dasar
				product.SatuanBesar = besar
				product.IsiPerBesar = isi
				product.HargaJualBesar = jualBesar // 🚀 4. SIMPAN HARGA GROSIR

				tx.Save(&product)
			} else {
				// BIKIN BARANG BARU
				newProduct := models.Product{
					StoreID: storeID, SKU: &sku, NamaProduk: nama, Kategori: kategori,
					HargaModal: modal, HargaJual: jual, Stok: stok, SatuanDasar: dasar,
					SatuanBesar: besar, IsiPerBesar: isi,
					HargaJualBesar: jualBesar, // 🚀 5. SIMPAN HARGA GROSIR
				}
				tx.Create(&newProduct)
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal impor: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil mengimpor " + strconv.Itoa(len(records)) + " produk"})
}

func (h *RetailHandler) GetDashboardReport(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Laporan keuangan cuma untuk Owner."})
		return
	}

	storeID := uint(c.MustGet("store_id").(float64))
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	now := time.Now()
	location := now.Location()

	start, _ := time.ParseInLocation("2006-01-02", startDateStr, location)
	if startDateStr == "" {
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	}

	end, _ := time.ParseInLocation("2006-01-02", endDateStr, location)
	if endDateStr == "" {
		end = start.Add(24 * time.Hour)
	} else {
		end = end.Add(24 * time.Hour)
	}

	// 🚀 STRUCT DI-UPDATE BUAT NAMPUNG DATA KLAIM & FINAL NETTO
	var report struct {
		TotalOmzet         float64 `json:"total_omzet"`
		TotalLaba          float64 `json:"total_laba"`
		JumlahTransaksi    int64   `json:"jumlah_transaksi"`
		TotalProdukTerjual float64 `json:"total_produk_terjual"`
		AvgTransaksi       float64 `json:"avg_transaksi"`
		TotalReturQty      float64 `json:"total_retur_qty"`
		TotalReturLoss     float64 `json:"total_retur_loss"`

		// Data SO Asli
		TotalSOQty  float64 `json:"total_so_qty"`
		TotalSOLoss float64 `json:"total_so_loss"`

		// Data Klaim Barang Nyempil
		TotalKlaimQty   float64 `json:"total_klaim_qty"`
		TotalKlaimValue float64 `json:"total_klaim_value"`

		// 🚀 HASIL AKHIR (FINAL) SETELAH DIKURANGI KLAIM
		NetSOQty  float64 `json:"net_so_qty"`
		NetSOLoss float64 `json:"net_so_loss"`
	}

	omzet, qty, _ := h.Repo.GetDashboardSummary(storeID, start, end)
	report.TotalOmzet = omzet
	report.TotalProdukTerjual = qty

	laba, _ := h.Repo.GetDashboardLaba(storeID, start, end)
	report.TotalLaba = laba

	db := h.Repo.GetDB()
	db.Model(&models.Transaction{}).Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, start, end).Count(&report.JumlahTransaksi)

	if report.JumlahTransaksi > 0 {
		report.AvgTransaksi = report.TotalOmzet / float64(report.JumlahTransaksi)
	}

	returQty, returLoss, _ := h.Repo.GetDashboardReturSummary(storeID, start, end)
	report.TotalReturQty = returQty
	report.TotalReturLoss = returLoss

	// 🚀 1. TARIK DATA SO AWAL
	soQty, soLoss, _ := h.Repo.GetDashboardSOSummary(storeID, start, end)
	report.TotalSOQty = soQty
	report.TotalSOLoss = soLoss

	// 🚀 2. TARIK DATA KLAIM YANG UDAH DI-APPROVE
	klaimQty, klaimValue, _ := h.Repo.GetDashboardKlaimSummary(storeID, start, end)
	report.TotalKlaimQty = klaimQty
	report.TotalKlaimValue = klaimValue

	// 🚀 3. HITUNG HASIL FINAL (NETTO RUGI)
	report.NetSOQty = soQty - klaimQty
	if report.NetSOQty < 0 {
		report.NetSOQty = 0
	} // Biar gak minus kalau anomali

	report.NetSOLoss = soLoss - klaimValue
	if report.NetSOLoss < 0 {
		report.NetSOLoss = 0
	} // Biar gak minus kalau anomali

	lowStock, _ := h.Repo.GetLowStockProducts(storeID, 10)

	type GrafikData struct {
		Tanggal   string  `json:"tanggal"`
		Omzet     float64 `json:"omzet"`
		Laba      float64 `json:"laba"`
		ReturLoss float64 `json:"retur_loss"`
	}
	var grafikPenjualan []GrafikData

	days := int(end.Sub(start).Hours() / 24)
	if days <= 0 {
		days = 1
	}
	if days > 31 {
		days = 31
	}

	for i := 0; i < days; i++ {
		tgl := start.AddDate(0, 0, i)
		tglEnd := tgl.Add(24 * time.Hour)

		dailyOmzet, dailyLaba, dailyReturLoss, _ := h.Repo.GetDailySalesReport(storeID, tgl, tglEnd)
		grafikPenjualan = append(grafikPenjualan, GrafikData{
			Tanggal:   tgl.Format("02 Jan"),
			Omzet:     dailyOmzet,
			Laba:      dailyLaba,
			ReturLoss: dailyReturLoss,
		})
	}

	bestSellers, _ := h.Repo.GetTopBestSellers(storeID, start, end)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"summary":          report,
			"grafik_penjualan": grafikPenjualan,
			"best_sellers":     bestSellers,
			"low_stock":        lowStock,
		},
	})
}

// 📅 SCHEDULE STRUCTURES & METHODS
type ScheduleItem struct {
	UserID    uint   `json:"user_id" binding:"required"`
	Tanggal   string `json:"tanggal" binding:"required"`
	ShiftType string `json:"shift_type" binding:"required"`
}

type BulkScheduleInput struct {
	Schedules []ScheduleItem `json:"schedules" binding:"required"`
}

func (h *RetailHandler) SaveSchedules(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	if val, ok := storeIDRaw.(float64); ok {
		storeID = uint(val)
	} else if val, ok := storeIDRaw.(uint); ok {
		storeID = val
	}

	var input BulkScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data jadwal tidak valid!"})
		return
	}

	db := h.Repo.GetDB()
	tx := db.Begin()

	for _, item := range input.Schedules {
		jamMasuk, jamPulang := "-", "-"
		if item.ShiftType == "Shift 1" {
			jamMasuk = "07:00"
			jamPulang = "15:00"
		} else if item.ShiftType == "Shift 2" {
			jamMasuk = "15:00"
			jamPulang = "23:00"
		} else if item.ShiftType == "Middle" {
			jamMasuk = "11:00"
			jamPulang = "19:00"
		}

		existing, err := h.Repo.GetScheduleByDate(tx, item.UserID, item.Tanggal)
		if err == nil {
			existing.ShiftType = item.ShiftType
			existing.JamMasukJadwal = jamMasuk
			existing.JamPulangJadwal = jamPulang
			if err := h.Repo.SaveScheduleTx(tx, existing); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui jadwal lama!"})
				return
			}
		} else {
			newSchedule := models.Schedule{
				StoreID:         storeID,
				UserID:          item.UserID,
				Tanggal:         item.Tanggal,
				ShiftType:       item.ShiftType,
				JamMasukJadwal:  jamMasuk,
				JamPulangJadwal: jamPulang,
			}
			if err := h.Repo.CreateScheduleTx(tx, &newSchedule); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jadwal baru!"})
				return
			}
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Jadwal mingguan berhasil disimpan! 🚀"})
}

func (h *RetailHandler) GetSchedules(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	if val, ok := storeIDRaw.(float64); ok {
		storeID = uint(val)
	} else if val, ok := storeIDRaw.(uint); ok {
		storeID = val
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	listJadwal, err := h.Repo.GetSchedulesRange(storeID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik data jadwal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": listJadwal})
}

// 💵 CASHIER SESSION STRUCTURES & METHODS
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
	session.EndTime = &now
	session.Status = "closed"

	h.Repo.SaveSession(session)

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

// 🛒 POS TRANSACTION LOGIC METHODS
type CartItem struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Kuantitas int     `json:"kuantitas" binding:"required,gt=0"`
	UomLabel  string  `json:"uom_label"`
	HargaUom  float64 `json:"harga_uom"`
}
type TransactionInput struct {
	Items         []CartItem `json:"items" binding:"required,gt=0"`
	NominalBayar  float64    `json:"nominal_bayar" binding:"required"`
	MetodeBayar   string     `json:"metode_bayar"`
	NoHPPelanggan string     `json:"no_hp_pelanggan"`
}

func (h *RetailHandler) CreateTransaction(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	userID := uint(c.MustGet("user_id").(float64))

	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format keranjang tidak sesuai!"})
		return
	}

	// 🚀 FUNGSI BANTUAN FORMAT RUPIAH OTOMATIS
	formatRupiah := func(amount float64) string {
		str := fmt.Sprintf("%.0f", amount)
		var result string
		for i, n := len(str)-1, 0; i >= 0; i-- {
			result = string(str[i]) + result
			n++
			if n%3 == 0 && i > 0 {
				result = "." + result
			}
		}
		return "Rp " + result
	}

	var savedTransaction models.Transaction

	db := h.Repo.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		activeSession, err := h.Repo.GetActiveSession(tx, userID, storeID)
		if err != nil {
			return fmt.Errorf("session kasir tidak ditemukan, silakan buka kasir dulu")
		}

		var store models.Store
		if err := tx.Select("id", "nama_toko", "business_type", "pajak_persen", "wa_token", "receipt_footer").First(&store, storeID).Error; err != nil {
			return fmt.Errorf("data toko tidak ditemukan")
		}

		tipeBisnis := "RETAIL"
		statusPesanan := "SELESAI"
		if store.BusinessType == "Jasa - Laundry" {
			tipeBisnis = "LAUNDRY"
			statusPesanan = "ANTRI"
		} else if store.BusinessType == "Kuliner - F&B" {
			tipeBisnis = "FNB"
			statusPesanan = "PROSES"
		}

		var subTotal float64
		var details []models.TransactionDetail
		rincianBarangWA := ""

		for _, item := range input.Items {
			product, err := h.Repo.GetProductByID(tx, item.ProductID, storeID)
			if err != nil {
				return err
			}

			if product.Stok < item.Kuantitas {
				return fmt.Errorf("Stok %s habis! Sisa: %d", product.NamaProduk, product.Stok)
			}

			// Kurangi stok berdasarkan total eceran (batang)
			product.Stok -= item.Kuantitas
			if err := h.Repo.SaveProduct(tx, product); err != nil {
				return err
			}

			// 🚀 LOGIKA KALKULASI DISPLAY
			itemSubTotal := float64(item.Kuantitas) * (product.HargaJual) // Fallback kalkulasi

			rincianDisplay := item.UomLabel
			hargaSatuanDisplay := item.HargaUom

			// Jika item.HargaUom dari Vue valid, gunakan untuk memastikan subtotal presisi
			if hargaSatuanDisplay > 0 {
				// Cari original Qty (misal "3" dari "3 BUNGKUS")
				qtyOriginal := 1
				fmt.Sscanf(rincianDisplay, "%d", &qtyOriginal)
				if qtyOriginal > 0 {
					itemSubTotal = float64(qtyOriginal) * hargaSatuanDisplay
				}
			}

			subTotal += itemSubTotal

			// Teks WA yang rapi dan detail
			rincianBarangWA += fmt.Sprintf("▪️ *%s*\n   %s x %s = *%s*\n",
				product.NamaProduk,
				rincianDisplay,
				formatRupiah(hargaSatuanDisplay),
				formatRupiah(itemSubTotal),
			)

			details = append(details, models.TransactionDetail{
				ProductID:   product.ID,
				HargaSatuan: hargaSatuanDisplay,
				Kuantitas:   item.Kuantitas, // Disimpan eceran untuk histori gudang
				SubTotal:    itemSubTotal,
				ItemType:    "PRODUCT",
				DetailNotes: rincianDisplay, // Disimpan format "3 BUNGKUS"
			})
		}

		pajak := (store.PajakPersen / 100.0) * subTotal
		rawTotal := subTotal + pajak
		roundedTotal := math.Round(rawTotal/100) * 100
		pembulatan := roundedTotal - rawTotal

		kembalian := input.NominalBayar - roundedTotal
		if kembalian < 0 {
			return fmt.Errorf("Uang pelanggan kurang!")
		}

		noInvoice := fmt.Sprintf("INV-%s", time.Now().Format("20060102150405"))

		savedTransaction = models.Transaction{
			SessionID:     activeSession.ID,
			StoreID:       storeID,
			UserID:        userID,
			NoInvoice:     noInvoice,
			SubTotal:      subTotal,
			Pajak:         pajak,
			Pembulatan:    pembulatan,
			TotalHarga:    roundedTotal,
			MetodeBayar:   input.MetodeBayar,
			StatusBayar:   "LUNAS",
			TipeBisnis:    tipeBisnis,
			StatusPesanan: statusPesanan,
			NominalBayar:  input.NominalBayar,
			Kembalian:     kembalian,
			Details:       details,
		}

		// 🚀 KIRIM WA DENGAN FORMAT RUPIAH YANG SEMPURNA
		if input.NoHPPelanggan != "" && store.WaToken != "" {
			pesanNota := fmt.Sprintf(
				`🏪 *%s*
========================
Halo Bosku! Terima kasih sudah berbelanja. Berikut rincian transaksi Anda:

🧾 *No. Invoice:* %s
📅 *Tanggal:* %s

*Rincian Pesanan:*
%s========================
💰 *Subtotal:* %s
⚖️ *Pajak/Biaya:* %s
========================
✅ *TOTAL BAYAR: %s*
💳 *Metode:* %s
💵 *Tunai:* %s
💸 *Kembali:* %s
========================
				%s`,
				store.NamaToko,
				noInvoice,
				time.Now().Format("02 Jan 2006, 15:04 WIB"),
				rincianBarangWA,
				formatRupiah(subTotal),
				formatRupiah(pajak),
				formatRupiah(roundedTotal),
				input.MetodeBayar,
				formatRupiah(input.NominalBayar),
				formatRupiah(kembalian),
				store.ReceiptFooter,
			)

			utils.SendWhatsAppFonnte(store.WaToken, input.NoHPPelanggan, pesanNota)
		}

		return h.Repo.CreateTransactionTx(tx, &savedTransaction)
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil! Struk siap dicetak.",
		"invoice": savedTransaction.NoInvoice,
		"tagihan": savedTransaction.TotalHarga,
		"kembali": savedTransaction.Kembalian,
	})
}

func (h *RetailHandler) GetTransactions(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	tanggal := c.Query("tanggal")
	if tanggal == "" {
		tanggal = time.Now().Format("2006-01-02")
	}

	parsedDate, err := time.ParseInLocation("2006-01-02", tanggal, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid"})
		return
	}

	startOfDay := parsedDate
	endOfDay := startOfDay.Add(24 * time.Hour)

	transactions, err := h.Repo.GetTransactionsByRange(storeID, startOfDay, endOfDay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik riwayat transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Riwayat transaksi berhasil ditarik!",
		"data":    transactions,
	})
}

// ==========================================
// ⚙️ PENGATURAN TOKO (STORE SETTINGS)
// ==========================================

func (h *RetailHandler) GetStoreSettings(c *gin.Context) {
	// Ambil store_id dari token JWT
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var store models.Store
	if err := src.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": store})
}

func (h *RetailHandler) UpdateStoreSettings(c *gin.Context) {
	// Proteksi: Hanya Owner & Manager yang bisa ubah settingan toko
	roleOwner, _ := c.Get("role")
	if roleOwner != "owner" && roleOwner != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Owner atau Manager yang bisa mengubah pengaturan toko!"})
		return
	}

	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))

	var store models.Store
	if err := src.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan!"})
		return
	}

	// 1. Update Data Teks Standar
	if v := c.PostForm("nama_toko"); v != "" {
		store.NamaToko = v
	}
	if v := c.PostForm("telepon"); v != "" {
		store.Telepon = utils.FormatPhoneNumber(v)
	}
	if v := c.PostForm("alamat"); v != "" {
		store.Alamat = v
	}
	if v := c.PostForm("provinsi"); v != "" {
		store.Provinsi = v
	}
	if v := c.PostForm("kota"); v != "" {
		store.Kota = v
	}
	if v := c.PostForm("kecamatan"); v != "" {
		store.Kecamatan = v
	}
	if v := c.PostForm("kelurahan"); v != "" {
		store.Kelurahan = v
	}
	if v := c.PostForm("kode_pos"); v != "" {
		store.KodePos = v
	}
	if v := c.PostForm("qris_name"); v != "" {
		store.QrisName = v
	}
	if v := c.PostForm("receipt_footer"); v != "" {
		store.ReceiptFooter = v
	}
	if v := c.PostForm("wa_token"); v != "" {
		store.WaToken = v
	}

	// 🚀 TAMBAHAN: Update Data Payment Gateway & Printer
	if v := c.PostForm("payment_type"); v != "" {
		store.PaymentType = v
	}
	if v := c.PostForm("midtrans_server_key"); v != "" {
		store.MidtransServerKey = v
	}
	if v := c.PostForm("midtrans_client_key"); v != "" {
		store.MidtransClientKey = v
	}
	if v := c.PostForm("printer_width"); v != "" {
		store.PrinterWidth = v
	}
	if v := c.PostForm("printer_type"); v != "" {
		store.PrinterType = v
	}

	// Toggle Pajak
	if v := c.PostForm("is_tax_active"); v != "" {
		store.IsTaxActive = (v == "true")
	}
	if v := c.PostForm("pajak_persen"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			store.PajakPersen = parsed
		}
	}

	// 2. 🚀 Update / Hapus Logo Struk
	if c.PostForm("delete_logo") == "true" {
		if store.LogoURL != "" {
			os.Remove("." + store.LogoURL) // Hapus file dari folder lokal server!
			store.LogoURL = ""             // Kosongkan URL di database
		}
	} else if file, err := c.FormFile("logo"); err == nil {
		newFileName := fmt.Sprintf("store_%d_logo_%d%s", storeID, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)
		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
			if store.LogoURL != "" {
				os.Remove("." + store.LogoURL)
			} // Hapus logo lama
			store.LogoURL = "/uploads/" + newFileName
		}
	}

	// 3. 🚀 Update / Hapus Barcode QRIS
	if c.PostForm("delete_qris") == "true" {
		if store.QrisImage != "" {
			os.Remove("." + store.QrisImage) // Hapus file fisik
			store.QrisImage = ""             // Kosongkan dari DB
		}
	} else if file, err := c.FormFile("qris"); err == nil {
		newFileName := fmt.Sprintf("store_%d_qris_%d%s", storeID, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)
		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
			if store.QrisImage != "" {
				os.Remove("." + store.QrisImage)
			} // Hapus qris lama
			store.QrisImage = "/uploads/" + newFileName
		}
	}

	// Simpan ke Database
	if err := src.DB.Save(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengaturan toko"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengaturan toko berhasil diperbarui!", "data": store})
}

// 🚀 STRUCT BUAT NANGKEP REQUEST DARI VUE
type UpgradeInput struct {
	PlanName string `json:"plan_name"`
	Price    int64  `json:"price"`
}

// 🚀 FUNGSI BIKIN TAGIHAN MIDTRANS
func (h *RetailHandler) CreateUpgradePayment(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))

	var input UpgradeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data paket tidak valid"})
		return
	}

	// 1. SETUP KUNCI RAHASIA MIDTRANS (Pakai Server Key Sandbox Lu)
	// Nanti ganti pakai Server Key asli dari dashboard Midtrans lu
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox // Default aman
	if os.Getenv("APP_ENV") == "production" {
		midtrans.Environment = midtrans.Production
	}

	// 2. BIKIN KERANJANG TAGIHAN
	orderID := fmt.Sprintf("UPGRADE-TOKO-%d-%s-%d", storeID, strings.ToUpper(input.PlanName), time.Now().Unix())

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: input.Price,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "SUB-" + strings.ToUpper(input.PlanName),
				Price: input.Price,
				Qty:   1,
				Name:  "Langganan Paket " + input.PlanName,
			},
		},
	}

	// 3. MINTA TOKEN KE MIDTRANS
	snapResp, err := snap.CreateTransaction(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungi Payment Gateway"})
		return
	}

	// 4. KASIH TOKENNYA KE VUE
	c.JSON(http.StatusOK, gin.H{"token": snapResp.Token, "order_id": orderID})
}

func (h *RetailHandler) MidtransWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid"})
		return
	}

	orderID, _ := payload["order_id"].(string)
	transactionStatus, _ := payload["transaction_status"].(string)

	if transactionStatus == "settlement" || transactionStatus == "capture" {
		parts := strings.Split(orderID, "-")

		// A. CEK UPGRADE PAKET
		if len(parts) >= 5 && parts[0] == "UPGRADE" {
			storeID := parts[2]
			planName := parts[3]
			endDate := time.Now().AddDate(0, 1, 0)
			db := h.Repo.GetDB()
			db.Exec("UPDATE stores SET subscription_status = ?, subscription_end = ?, subscription_plan = ? WHERE id = ?",
				"active", endDate, strings.ToLower(planName), storeID)

			// B. CEK TRANSAKSI POS (KASIR)
		} else if len(parts) >= 2 && parts[0] == "POS" {
			db := h.Repo.GetDB()
			err := db.Exec("UPDATE transactions SET status_bayar = ? WHERE no_invoice = ?",
				"LUNAS", orderID).Error

			if err != nil {
				fmt.Println("❌ GAGAL UPDATE TRANSAKSI:", err)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ==============================================================
// 🚀 MIDTRANS TRANSAKSI KASIR (UANG MASUK KE REKENING TENANT/TOKO)
// ==============================================================

// Struct buat nangkep total belanja dari Vue
type PosMidtransReq struct {
	Total float64 `json:"total" binding:"required"`
}

func (h *RetailHandler) CreatePosMidtransOrder(c *gin.Context) {
	storeIDRaw, exists := c.Get("store_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}

	var storeID uint
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
	}

	var input PosMidtransReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data nominal tidak valid"})
		return
	}

	var store models.Store
	if err := src.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan!"})
		return
	}

	if store.PaymentType != "midtrans" || store.MidtransServerKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Toko belum mengatur Midtrans Server Key!"})
		return
	}

	// 1. INIT CLIENT MIDTRANS
	var s snap.Client

	// 🚀 FIX: PAKSA JADI SANDBOX (HAPUS LOGIKA DETEKSI SB-)
	env := midtrans.Sandbox
	if os.Getenv("APP_ENV") == "production" {
		env = midtrans.Production
	}

	s.New(store.MidtransServerKey, env)

	orderID := fmt.Sprintf("POS-STR%d-%d", storeID, time.Now().Unix())

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(input.Total),
		},
	}

	// TEMBAK API MIDTRANS
	snapResp, err := s.CreateTransaction(req)
	if err != nil {
		fmt.Println("❌ ERROR MIDTRANS:", err.GetMessage())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.GetMessage()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": snapResp.Token, "order_id": orderID})
}
