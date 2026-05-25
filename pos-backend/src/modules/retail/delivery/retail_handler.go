package delivery

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
	"pos-backend/models"
	"golang.org/x/crypto/bcrypt"
	"math"

	"pos-backend/src/modules/retail/domain"
	"pos-backend/src/modules/retail/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RetailHandler struct {
	Repo repository.RetailRepository
}

func NewRetailHandler(repo repository.RetailRepository) *RetailHandler {
	return &RetailHandler{Repo: repo}
}

// 🚀 STOCK OPNAME HANDLERS
type StockOpnameInput struct {
	Notes string `json:"notes"`
	Items []struct {
		ProductID uint `json:"product_id"`
		ActualQty int  `json:"actual_qty"`
	} `json:"items"`
}

func (h *RetailHandler) CreateStockOpname(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	userID := uint(c.MustGet("user_id").(float64))

	var input StockOpnameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid!"})
		return
	}

	if len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Daftar barang opname tidak boleh kosong!"})
		return
	}

	db := h.Repo.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		so := domain.StockOpname{
			StoreID:   storeID,
			UserID:    userID,
			Notes:     input.Notes,
			CreatedAt: time.Now(),
		}

		if err := h.Repo.CreateStockOpname(tx, &so); err != nil {
			return err
		}

		for _, item := range input.Items {
			product, err := h.Repo.GetProductByID(tx, item.ProductID, storeID)
			if err != nil {
				return err 
			}

			selisih := item.ActualQty - product.Stok

			detail := domain.StockOpnameDetail{
				OpnameID:  so.ID,
				ProductID: item.ProductID,
				SystemQty: product.Stok,
				ActualQty: item.ActualQty,
				Selisih:   selisih,
			}

			if err := h.Repo.CreateStockOpnameDetail(tx, &detail); err != nil {
				return err
			}

			product.Stok = item.ActualQty
			if err := h.Repo.SaveProduct(tx, product); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses Stock Opname: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock Opname berhasil disimpan. Stok master telah diperbarui!"})
}

func (h *RetailHandler) GetStockOpnameHistory(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	history, err := h.Repo.GetStockOpnameHistory(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data riwayat opname"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": history})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	db := h.Repo.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		purchase := domain.Purchase{
			StoreID:      storeID,
			UserID:       userID,
			SupplierName: input.SupplierName,
			NoFaktur:     input.NoFaktur,
			TotalItem:    len(input.Items),
		}

		if err := h.Repo.CreatePurchase(tx, &purchase); err != nil {
			return err
		}

		for _, item := range input.Items {
			product, err := h.Repo.GetProductByID(tx, item.ProductID, storeID)
			if err != nil {
				return err
			}

			product.Stok += item.QtyMasuk
			product.HargaModal = item.HargaModal

			if err := h.Repo.SaveProduct(tx, product); err != nil {
				return err
			}

			detail := domain.PurchaseDetail{
				PurchaseID: purchase.ID,
				ProductID:  item.ProductID,
				QtyMasuk:   item.QtyMasuk,
				HargaModal: item.HargaModal,
			}
			if err := h.Repo.CreatePurchaseDetail(tx, &detail); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal proses LPB: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "LPB Berhasil! Stok barang sudah bertambah."})
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
	if val, ok := storeIDRaw.(float64); ok { storeID = uint(val) } else if val, ok := storeIDRaw.(uint); ok { storeID = val }

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
	if val, ok := storeIDRaw.(float64); ok { storeID = uint(val) } else if val, ok := storeIDRaw.(uint); ok { storeID = val }

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
		if logData.User.NIK != nil { nik = *logData.User.NIK }

		jamPulang := logData.JamPulang
		if jamPulang == "" { jamPulang = "Belum Pulang" }

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

	if name == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Password wajib diisi!"})
		return
	}
	if inputRole == "" { inputRole = "kasir" }

	currentYear := time.Now().Format("2006")
	var newNIK string

	lastEmployee, err := h.Repo.GetLastEmployeeNIK(storeID, currentYear)
	if err != nil {
		newNIK = currentYear + "0001"
	} else {
		lastNIK := *lastEmployee.NIK
		lastSequence, _ := strconv.Atoi(lastNIK[4:])
		newNIK = fmt.Sprintf("%s%04d", currentYear, lastSequence+1)
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
		NoHP:         noHP,
		FotoURL:      fotoURL,
	}

	if err := h.Repo.CreateEmployee(&employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan ke database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Karyawan baru berhasil didaftarkan! 🤝",
		"data": gin.H{"nama": employee.Name, "nik": newNIK, "jabatan": employee.Role},
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
	employee.NoHP = c.PostForm("no_hp")
	
	if newRole := c.PostForm("role"); newRole != "" { employee.Role = newRole }
	if password := c.PostForm("password"); password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		employee.Password = string(hashed)
	}

	nikClean := "karyawan"
	if employee.NIK != nil { nikClean = *employee.NIK }

	if file, err := c.FormFile("foto"); err == nil {
		newFileName := fmt.Sprintf("%s_%d%s", nikClean, time.Now().Unix(), filepath.Ext(file.Filename))
		uploadPath := filepath.Join("uploads", newFileName)
		if err := c.SaveUploadedFile(file, uploadPath); err == nil {
			if employee.FotoURL != "" { os.Remove("." + employee.FotoURL) }
			employee.FotoURL = "/uploads/" + newFileName
		}
	}

	if bioFile, errBio := c.FormFile("biometric_file"); errBio == nil {
		newBioName := fmt.Sprintf("%s_bio_%d%s", nikClean, time.Now().Unix(), filepath.Ext(bioFile.Filename))
		uploadBioPath := filepath.Join("uploads", newBioName)
		if err := c.SaveUploadedFile(bioFile, uploadBioPath); err == nil {
			if employee.BiometricURL != "" { os.Remove("." + employee.BiometricURL) }
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
	SKU         *string `form:"sku"`
	NamaProduk  string  `form:"nama_produk" binding:"required"`
	Kategori    string  `form:"kategori"`
	HargaModal  float64 `form:"harga_modal"`
	HargaJual   float64 `form:"harga_jual" binding:"required"`
	Stok        int     `form:"stok"`
	SatuanDasar string  `form:"satuan_dasar"`
	SatuanBesar string  `form:"satuan_besar"`
	IsiPerBesar int     `form:"isi_per_besar"`
}

func (h *RetailHandler) CreateProduct(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))

	var input ProductInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Periksa kembali isian form Anda: " + err.Error()})
		return
	}
	if input.SatuanDasar == "" { input.SatuanDasar = "PCS" }

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
		StoreID:     storeID,
		SKU:         input.SKU,
		NamaProduk:  input.NamaProduk,
		Kategori:    input.Kategori,
		HargaModal:  input.HargaModal,
		HargaJual:   input.HargaJual,
		Stok:        input.Stok,
		Gambar:      imagePath,
		SatuanDasar: input.SatuanDasar,
		SatuanBesar: input.SatuanBesar,
		IsiPerBesar: input.IsiPerBesar,
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
		if pageStr == "" { pageStr = "1" }
		if limitStr == "" { limitStr = "10" }
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
	if sku := c.PostForm("sku"); sku != "" { product.SKU = &sku } else { product.SKU = nil }
	product.Kategori = c.PostForm("kategori")

	if hargaModal, err := strconv.ParseFloat(c.PostForm("harga_modal"), 64); err == nil { product.HargaModal = hargaModal }
	if hargaJual, err := strconv.ParseFloat(c.PostForm("harga_jual"), 64); err == nil { product.HargaJual = hargaJual }
	if stok, err := strconv.Atoi(c.PostForm("stok")); err == nil { product.Stok = stok }
	if satuanDasar := c.PostForm("satuan_dasar"); satuanDasar != "" { product.SatuanDasar = satuanDasar }
	product.SatuanBesar = c.PostForm("satuan_besar")
	if isiPerBesar, err := strconv.Atoi(c.PostForm("isi_per_besar")); err == nil { product.IsiPerBesar = isiPerBesar } else { product.IsiPerBesar = 0 }

	file, errFile := c.FormFile("gambar")
	if errFile == nil {
		folderPath := "uploads/products"
		os.MkdirAll(folderPath, os.ModePerm)
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		savePath := filepath.Join(folderPath, fileName)

		if err := c.SaveUploadedFile(file, savePath); err == nil {
			if product.Gambar != "" { os.Remove("." + product.Gambar) }
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

	if product.Gambar != "" { os.Remove("." + product.Gambar) }
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
	w.Write([]string{"SKU", "Nama Produk", "Kategori", "Harga Modal", "Harga Jual", "Stok", "Satuan Dasar", "Satuan Besar", "Isi Per Besar"})

	for _, p := range products {
		sku := ""
		if p.SKU != nil { sku = *p.SKU }
		w.Write([]string{
			sku, p.NamaProduk, p.Kategori,
			fmt.Sprintf("%.0f", p.HargaModal), fmt.Sprintf("%.0f", p.HargaJual),
			fmt.Sprintf("%d", p.Stok), p.SatuanDasar, p.SatuanBesar, fmt.Sprintf("%d", p.IsiPerBesar),
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
	_, _ = reader.Read() 

	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca isi CSV"})
		return
	}

	db := h.Repo.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, row := range records {
			sku := row[0]; nama := row[1]; kategori := row[2]
			modal, _ := strconv.ParseFloat(row[3], 64)
			jual, _ := strconv.ParseFloat(row[4], 64)
			stok, _ := strconv.Atoi(row[5])
			dasar := row[6]; besar := row[7]; isi, _ := strconv.Atoi(row[8])

			if nama == "" { continue }

			var product models.Product
			res := tx.Where("sku = ? AND store_id = ?", sku, storeID).First(&product)

			if res.Error == nil {
				product.NamaProduk = nama; product.Kategori = kategori; product.HargaModal = modal
				product.HargaJual = jual; product.Stok = stok; product.SatuanDasar = dasar
				product.SatuanBesar = besar; product.IsiPerBesar = isi
				tx.Save(&product)
			} else {
				newProduct := models.Product{
					StoreID: storeID, SKU: &sku, NamaProduk: nama, Kategori: kategori,
					HargaModal: modal, HargaJual: jual, Stok: stok, SatuanDasar: dasar,
					SatuanBesar: besar, IsiPerBesar: isi,
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
	if startDateStr == "" { start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location) }

	end, _ := time.ParseInLocation("2006-01-02", endDateStr, location)
	if endDateStr == "" { end = start.Add(24 * time.Hour) } else { end = end.Add(24 * time.Hour) }

	var report struct {
		TotalOmzet         float64 `json:"total_omzet"`
		TotalLaba          float64 `json:"total_laba"`
		JumlahTransaksi    int64   `json:"jumlah_transaksi"`
		TotalProdukTerjual float64 `json:"total_produk_terjual"`
		AvgTransaksi       float64 `json:"avg_transaksi"`
		TotalReturQty      float64 `json:"total_retur_qty"`
		TotalReturLoss     float64 `json:"total_retur_loss"`
		TotalSOQty         float64 `json:"total_so_qty"`
		TotalSOLoss        float64 `json:"total_so_loss"`
	}

	omzet, qty, _ := h.Repo.GetDashboardSummary(storeID, start, end)
	report.TotalOmzet = omzet
	report.TotalProdukTerjual = qty

	laba, _ := h.Repo.GetDashboardLaba(storeID, start, end)
	report.TotalLaba = laba

	db := h.Repo.GetDB()
	db.Model(&models.Transaction{}).Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, start, end).Count(&report.JumlahTransaksi)

	if report.JumlahTransaksi > 0 { report.AvgTransaksi = report.TotalOmzet / float64(report.JumlahTransaksi) }

	returQty, returLoss, _ := h.Repo.GetDashboardReturSummary(storeID, start, end)
	report.TotalReturQty = returQty
	report.TotalReturLoss = returLoss

	soQty, soLoss, _ := h.Repo.GetDashboardSOSummary(storeID, start, end)
	report.TotalSOQty = soQty
	report.TotalSOLoss = soLoss

	lowStock, _ := h.Repo.GetLowStockProducts(storeID, 10)

	type GrafikData struct {
		Tanggal   string  `json:"tanggal"`
		Omzet     float64 `json:"omzet"`
		Laba      float64 `json:"laba"`
		ReturLoss float64 `json:"retur_loss"`
	}
	var grafikPenjualan []GrafikData

	days := int(end.Sub(start).Hours() / 24)
	if days <= 0 { days = 1 }
	if days > 31 { days = 31 }

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
	if val, ok := storeIDRaw.(float64); ok { storeID = uint(val) } else if val, ok := storeIDRaw.(uint); ok { storeID = val }

	var input BulkScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data jadwal tidak valid!"})
		return
	}

	db := h.Repo.GetDB()
	tx := db.Begin()

	for _, item := range input.Schedules {
		jamMasuk, jamPulang := "-", "-"
		if item.ShiftType == "Shift 1" { jamMasuk = "07:00"; jamPulang = "15:00" } else if item.ShiftType == "Shift 2" { jamMasuk = "15:00"; jamPulang = "23:00" } else if item.ShiftType == "Middle" { jamMasuk = "11:00"; jamPulang = "19:00" }

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
	if val, ok := storeIDRaw.(float64); ok { storeID = uint(val) } else if val, ok := storeIDRaw.(uint); ok { storeID = val }

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
	if exists { userRole = userRoleRaw.(string) }

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

	c.JSON(http.StatusOK, gin.H{"has_session": true, "session": session})
}

func (h *RetailHandler) CloseSession(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, _ := strconv.Atoi(sessionIDStr)
	
	var input struct { TotalAktual float64 `json:"total_aktual"` }
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
	salesQRIS, _ := h.Repo.GetSalesMethodSummary(sessionIDStr, "QRIS")
	salesBCA, _ := h.Repo.GetSalesMethodSummary(sessionIDStr, "BCA")
	salesNonTunai := salesQRIS + salesBCA

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
		"start_time":       session.StartTime.In(loc).Format("02.01.06 15:04"),
		"end_time":         session.EndTime.In(loc).Format("02.01.06 15:04"),
		"sales_gross":      salesGross,
		"total_tax":        totalTax,
		"net_sales":        netSales,
		"modal_awal":        session.ModalAwal,
		"sales_cash":       salesCash,
		"sales_qris":       salesQRIS,
		"sales_bca":        salesBCA,
		"sales_non_tunai":  salesNonTunai,
		"total_expected":   totalExpected,
		"total_actual":     input.TotalAktual,
		"selisih":          selisih,
	})
}

// 🛒 POS TRANSACTION LOGIC METHODS
type CartItem struct {
	ProductID uint `json:"product_id" binding:"required"`
	Kuantitas int  `json:"kuantitas" binding:"required,gt=0"`
}

type TransactionInput struct {
	Items        []CartItem `json:"items" binding:"required,gt=0"`
	NominalBayar float64    `json:"nominal_bayar" binding:"required"`
	MetodeBayar  string     `json:"metode_bayar"`
}

func (h *RetailHandler) CreateTransaction(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	userID := uint(c.MustGet("user_id").(float64))

	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format keranjang tidak sesuai!"})
		return
	}

	var savedTransaction models.Transaction

	db := h.Repo.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		activeSession, err := h.Repo.GetActiveSession(tx, userID, storeID)
		if err != nil {
			return fmt.Errorf("session kasir tidak ditemukan, silakan buka kasir dulu")
		}

		store, err := h.Repo.GetStoreByIDSimple(tx, storeID)
		if err != nil { return err }

		// 🚀 1. DETEKSI OTOMATIS TIPE BISNIS & STATUSNYA DARI CORE STORE
		tipeBisnis := "RETAIL"
		statusPesanan := "SELESAI" // Default Retail langsung lunas dibawa pulang

		if store.BusinessType == "Jasa - Laundry" {
			tipeBisnis = "LAUNDRY"
			statusPesanan = "ANTRI"   // Kalau Laundry, statusnya masuk antrean cuci
		} else if store.BusinessType == "Kuliner - F&B" {
			tipeBisnis = "FNB"
			statusPesanan = "PROSES"  // Kalau F&B, statusnya masuk antrean dapur
		}

		var subTotal float64
		var details []models.TransactionDetail

		for _, item := range input.Items {
			product, err := h.Repo.GetProductByID(tx, item.ProductID, storeID)
			if err != nil {
				return fmt.Errorf("barang dengan ID %d tidak ditemukan", item.ProductID)
			}

			if product.Stok < item.Kuantitas {
				return fmt.Errorf("Stok %s habis/kurang! Sisa Stok: %d", product.NamaProduk, product.Stok)
			}

			// Potong stok master inventory global
			product.Stok -= item.Kuantitas
			if err := h.Repo.SaveProduct(tx, product); err != nil { return err }

			itemSubTotal := product.HargaJual * float64(item.Kuantitas)
			subTotal += itemSubTotal

			// 🚀 2. MAPPING DETAIL ITEM SESUAI FORMAT ARSITEKTUR GLOBAL
			details = append(details, models.TransactionDetail{
				ProductID:   product.ID,
				HargaSatuan: product.HargaJual,
				Kuantitas:   item.Kuantitas,
				SubTotal:    itemSubTotal,
				ItemType:    "PRODUCT", // Default core retail
				DetailNotes: "Transaksi Retail Toko", // Pengganti kolom parfum
			})
		}

		pajak := (store.PajakPersen / 100.0) * subTotal
		rawTotal := subTotal + pajak

		roundedTotal := math.Round(rawTotal/100) * 100
		pembulatan := roundedTotal - rawTotal

		kembalian := input.NominalBayar - roundedTotal
		if kembalian < 0 {
			return fmt.Errorf("Uang pelanggan kurang Rp %.0f !", math.Abs(kembalian))
		}

		noInvoice := fmt.Sprintf("INV-%s", time.Now().Format("20060102150405"))

		// 🚀 3. SET KEPALA STRUK GLOBAL (HEADER)
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
			TipeBisnis:    tipeBisnis,    // 🟢 Mengikuti jenis toko ("RETAIL")
			StatusPesanan: statusPesanan, // 🟢 Mengikuti jenis toko ("SELESAI")
			NominalBayar:  input.NominalBayar,
			Kembalian:     kembalian,
			Details:       details,
		}

		return h.Repo.CreateTransactionTx(tx, &savedTransaction)
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil! 💸 Struk siap dicetak.",
		"invoice": savedTransaction.NoInvoice,
		"tagihan": savedTransaction.TotalHarga,
		"kembali": savedTransaction.Kembalian,
		"data":    savedTransaction,
	})
}

func (h *RetailHandler) GetTransactions(c *gin.Context) {
	storeID := uint(c.MustGet("store_id").(float64))
	tanggal := c.Query("tanggal")
	if tanggal == "" { tanggal = time.Now().Format("2006-01-02") }

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