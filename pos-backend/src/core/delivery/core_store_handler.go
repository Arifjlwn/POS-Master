package delivery

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"pos-backend/models"
	"pos-backend/src/core/repository"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
)

type StoreCoreHandler struct {
	Repo repository.CoreRepo
}

func NewStoreCoreHandler(repo repository.CoreRepo) *StoreCoreHandler {
	return &StoreCoreHandler{Repo: repo}
}

type StoreSettingsResponse struct {
	ID                 uint       `json:"id"`
	PublicID           string     `json:"public_id"`
	NamaToko           string     `json:"nama_toko"`
	Telepon            string     `json:"telepon"`
	BusinessType       string     `json:"business_type"`
	Industry           string     `json:"industry"`
	SubscriptionPlan   string     `json:"subscription_plan"`
	SubscriptionStatus string     `json:"subscription_status"`
	SubscriptionEnd    *time.Time `json:"subscription_end"`
	FiturAktif         string     `json:"fitur_aktif"`
	QuotaTerminal      int        `json:"quota_terminal"`
	Alamat             string     `json:"alamat"`
	Provinsi           string     `json:"provinsi"`
	Kota               string     `json:"kota"`
	Kecamatan          string     `json:"kecamatan"`
	Kelurahan          string     `json:"kelurahan"`
	KodePos            string     `json:"kode_pos"`
	Latitude           float64    `json:"latitude"`
	Longitude          float64    `json:"longitude"`
	LogoURL            string     `json:"logo_url"`
	IsTaxActive        bool       `json:"is_tax_active"`
	PajakPersen        float64    `json:"pajak_persen"`
	ReceiptFooter      string     `json:"receipt_footer"`
	PaymentType        string     `json:"payment_type"`
	QrisImage          string     `json:"qris_image"`
	QrisName           string     `json:"qris_name"`
	PrinterWidth       string     `json:"printer_width"`
	PrinterType        string     `json:"printer_type"`
	WaToken            string     `json:"wa_token"`
	MidtransServerKey  string     `json:"midtrans_server_key"`
	MidtransClientKey  string     `json:"midtrans_client_key"`
}

func (h *StoreCoreHandler) GetStoreSettings(c *gin.Context) {
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

	var store models.Store
	db := h.Repo.GetDB()
	if err := db.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan!"})
		return
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// 🚀 SECURITY FIX DOUBLE URL GAIB BRAY!
	qrisURL := store.QrisImage
	if qrisURL != "" && !strings.HasPrefix(qrisURL, "http://") && !strings.HasPrefix(qrisURL, "https://") {
		qrisURL = baseURL + "/" + qrisURL
	}

	logoURL := store.LogoURL
	if logoURL != "" && !strings.HasPrefix(logoURL, "http://") && !strings.HasPrefix(logoURL, "https://") {
		logoURL = baseURL + "/" + logoURL
	}

	resp := StoreSettingsResponse{
		ID:                 store.ID,
		PublicID:           store.PublicID,
		NamaToko:           store.NamaToko,
		Telepon:            store.Telepon,
		BusinessType:       store.BusinessType,
		Industry:           store.Industry,
		SubscriptionPlan:   store.SubscriptionPlan,
		SubscriptionStatus: store.SubscriptionStatus,
		SubscriptionEnd:    store.SubscriptionEnd,
		FiturAktif:         store.FiturAktif,
		QuotaTerminal:      store.QuotaTerminal,
		Alamat:             store.Alamat,
		Provinsi:           store.Provinsi,
		Kota:               store.Kota,
		Kecamatan:          store.Kecamatan,
		Kelurahan:          store.Kelurahan,
		KodePos:            store.KodePos,
		Latitude:           store.Latitude,
		Longitude:          store.Longitude,
		LogoURL:            logoURL,
		IsTaxActive:        store.IsTaxActive,
		PajakPersen:        store.PajakPersen,
		ReceiptFooter:      store.ReceiptFooter,
		PaymentType:        store.PaymentType,
		QrisImage:          qrisURL,
		QrisName:           store.QrisName,
		PrinterWidth:       store.PrinterWidth,
		PrinterType:        store.PrinterType,
	}

	if store.WaToken != "" {
		resp.WaToken = "HAS_TOKEN_HIDDEN_BY_SYSTEM"
	}
	if store.MidtransServerKey != "" {
		resp.MidtransServerKey = "HAS_KEY_HIDDEN_BY_SYSTEM"
	}
	if store.MidtransClientKey != "" {
		resp.MidtransClientKey = "HAS_KEY_HIDDEN_BY_SYSTEM"
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *StoreCoreHandler) UpdateStoreSettings(c *gin.Context) {
	roleOwner, _ := c.Get("role")
	if roleOwner != "owner" && roleOwner != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Owner atau Manager yang bisa mengubah pengaturan toko!"})
		return
	}

	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
	}

	db := h.Repo.GetDB()
	var store models.Store
	if err := db.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan!"})
		return
	}

	if v, ok := c.GetPostForm("nama_toko"); ok {
		store.NamaToko = v
	}
	if v, ok := c.GetPostForm("telepon"); ok {
		store.Telepon = utils.FormatPhoneNumber(v)
	}
	if v, ok := c.GetPostForm("alamat"); ok {
		store.Alamat = v
	}
	if v, ok := c.GetPostForm("provinsi"); ok {
		store.Provinsi = v
	}
	if v, ok := c.GetPostForm("kota"); ok {
		store.Kota = v
	}
	if v, ok := c.GetPostForm("kecamatan"); ok {
		store.Kecamatan = v
	}
	if v, ok := c.GetPostForm("kelurahan"); ok {
		store.Kelurahan = v
	}
	if v, ok := c.GetPostForm("kode_pos"); ok {
		store.KodePos = v
	}
	if v, ok := c.GetPostForm("qris_name"); ok {
		store.QrisName = v
	}
	if v, ok := c.GetPostForm("receipt_footer"); ok {
		store.ReceiptFooter = v
	}
	if v, ok := c.GetPostForm("payment_type"); ok {
		store.PaymentType = v
	}
	if v, ok := c.GetPostForm("printer_width"); ok {
		store.PrinterWidth = v
	}
	if v, ok := c.GetPostForm("printer_type"); ok {
		store.PrinterType = v
	}

	if v, ok := c.GetPostForm("latitude"); ok {
		if parsedLat, err := strconv.ParseFloat(v, 64); err == nil {
			store.Latitude = parsedLat
		}
	}
	if v, ok := c.GetPostForm("longitude"); ok {
		if parsedLng, err := strconv.ParseFloat(v, 64); err == nil {
			store.Longitude = parsedLng
		}
	}

	if v, ok := c.GetPostForm("wa_token"); ok && v != "HAS_TOKEN_HIDDEN_BY_SYSTEM" {
		store.WaToken = strings.TrimSpace(v)
	}
	if v, ok := c.GetPostForm("midtrans_server_key"); ok && v != "HAS_KEY_HIDDEN_BY_SYSTEM" {
		store.MidtransServerKey = strings.TrimSpace(v)
	}
	if v, ok := c.GetPostForm("midtrans_client_key"); ok && v != "HAS_KEY_HIDDEN_BY_SYSTEM" {
		store.MidtransClientKey = strings.TrimSpace(v)
	}

	if v, ok := c.GetPostForm("is_tax_active"); ok {
		store.IsTaxActive = (v == "true")
	}
	if v, ok := c.GetPostForm("pajak_persen"); ok {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			store.PajakPersen = parsed
		}
	}

	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")

	if c.PostForm("delete_logo") == "true" {
		store.LogoURL = ""
	} else if file, err := c.FormFile("logo"); err == nil {
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran logo maksimal 5 MB"})
			return
		}
		fileSrc, _ := file.Open()
		defer fileSrc.Close()
		remotePath := fmt.Sprintf("stores/%s/logo", store.PublicID)
		urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, file.Header.Get("Content-Type"), bucketName, remotePath)
		if errUpload == nil {
			store.LogoURL = urlResult
		}
	}

	if c.PostForm("delete_qris") == "true" {
		store.QrisImage = ""
	} else if file, err := c.FormFile("qris"); err == nil {
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran QRIS maksimal 5 MB"})
			return
		}
		fileSrc, _ := file.Open()
		defer fileSrc.Close()
		remotePath := fmt.Sprintf("stores/%s/qris", store.PublicID)
		urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, file.Header.Get("Content-Type"), bucketName, remotePath)
		if errUpload == nil {
			store.QrisImage = urlResult
		}
	}

	if err := db.Save(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengaturan toko"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pengaturan toko berhasil diperbarui!", "data": store})
}

func (h *StoreCoreHandler) TestWhatsApp(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var store models.Store
	db := h.Repo.GetDB()
	if err := db.First(&store, storeIDRaw).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan"})
		return
	}

	var body struct {
		Target string `json:"target"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor target wajib diisi"})
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	payload := strings.NewReader(fmt.Sprintf("target=%s&message=ARZURA POS: Integrasi berhasil!", body.Target))
	req, _ := http.NewRequest("POST", "https://api.fonnte.com/send", payload)
	req.Header.Add("Authorization", store.WaToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungi Fonnte"})
		return
	}
	defer resp.Body.Close()
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Pesan uji coba terkirim!"})
}
