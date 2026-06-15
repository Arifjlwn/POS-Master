package delivery

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"
	"pos-backend/src/modules/jasalayanan/laundry/repository"
	"pos-backend/src/modules/jasalayanan/laundry/usecase"
	"github.com/gin-gonic/gin"
)

type LaundryReportHandler struct {
	usecase usecase.LaundryUseCase
	Repo    repository.LaundryRepository
}

func NewLaundryReportHandler(u usecase.LaundryUseCase, r repository.LaundryRepository) *LaundryReportHandler {
	return &LaundryReportHandler{usecase: u, Repo: r}
}

type UpdateSettingInput struct {
	NamaToko      string `json:"nama_toko" binding:"required"`
	Telepon       string `json:"telepon" binding:"required"`
	Alamat        string `json:"alamat" binding:"required"`
	PaymentType   string `json:"payment_type" binding:"required"`
	QrisBase64    string `json:"qris_base64"`
	ReceiptFooter string `json:"receipt_footer"`
}

func (h *LaundryReportHandler) CariPelanggan(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	keyword := c.Query("q")
	customers, err := h.Repo.SearchCustomers(storeID, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari data pelanggan"})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (h *LaundryReportHandler) GetSettingToko(c *gin.Context) {
    storeIDRaw, _ := c.Get("store_id")
    storeID := uint(storeIDRaw.(float64))
    
    store, err := h.Repo.GetStoreByID(storeID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Data toko tidak ditemukan"})
        return
    }
    
    baseURL := os.Getenv("BASE_URL")
    if baseURL == "" { 
        baseURL = "http://localhost:8080" 
    }
    if store.QrisImage != "" { 
        store.QrisImage = baseURL + "/" + store.QrisImage 
    }
    
    // 🚀 UPGRADE SINKRONISASI BRAY: Masukin semua field billing SaaS ke JSON Response
    c.JSON(http.StatusOK, gin.H{
        "id":                  store.ID, 
        "nama_toko":           store.NamaToko, 
        "business_type":       store.BusinessType, 
        "industry":            store.Industry, // Ikut di-expose bray
        "telepon":             store.Telepon, 
        "alamat":              store.Alamat, 
        "payment_type":        store.PaymentType, 
        "qris_image":          store.QrisImage, 
        "receipt_footer":      store.ReceiptFooter,
        "quota_terminal":      store.QuotaTerminal,
        
        // 🔒 Kunci Utama Paywall Billing Lu Rif!
        "subscription_plan":   store.SubscriptionPlan,
        "subscription_status": store.SubscriptionStatus,
        "subscription_end":    store.SubscriptionEnd, // Nanti dikirim bentuk format ISO-Time/Null otomatis bray
    })
}

func (h *LaundryReportHandler) UpdateSettingToko(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	var input UpdateSettingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data pengaturan tidak valid"})
		return
	}
	db := h.Repo.GetDB()
	tx := db.Begin()
	store, err := h.Repo.GetStoreByID(storeID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan"})
		return
	}
	if input.PaymentType == "PRIBADI" && input.QrisBase64 != "" && !strings.HasPrefix(input.QrisBase64, "http") {
		parts := strings.Split(input.QrisBase64, ",")
		pureBase64 := parts[0]
		if len(parts) > 1 { pureBase64 = parts[1] }
		decodedData, _ := base64.StdEncoding.DecodeString(pureBase64)
		os.MkdirAll("public/uploads/qris_toko", os.ModePerm)
		qrisPath := "public/uploads/qris_toko/store_qris.jpg"
		os.WriteFile(qrisPath, decodedData, 0644)
		store.QrisImage = qrisPath
	}
	store.NamaToko = input.NamaToko; store.Telepon = input.Telepon; store.Alamat = input.Alamat; store.PaymentType = input.PaymentType; store.ReceiptFooter = input.ReceiptFooter
	if err := h.Repo.UpdateStoreTx(tx, store); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui profil toko"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "Pengaturan toko berhasil diperbarui!"})
}

func (h *LaundryReportHandler) AmbilDataTracking(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	
	results, err := h.Repo.GetTrackingCucian(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tracking"})
		return // 🌟 Kalau di sini return-nya bener bray murni buat nge-stop eksekusi fungsi ke bawah!
	}
	
	// 🚀 FIX: Buang kata "return"-nya bray! Biar murni nembak response Gin
	c.JSON(http.StatusOK, results)
}