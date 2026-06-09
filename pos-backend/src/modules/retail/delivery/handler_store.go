package delivery

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"pos-backend/models"
	src "pos-backend/src/core/config"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// Struct khusus untuk mengamankan data sensitif agar tidak bocor ke kasir/staf toko
// Taruh ini di atas fungsi GetStoreSettings di handle_store.go
type StoreSettingsResponse struct {
	ID           uint   `json:"id"`
	PublicID     string `json:"public_id"`
	NamaToko     string `json:"nama_toko"`
	Telepon      string `json:"telepon"`
	BusinessType string `json:"business_type"`
	Industry     string `json:"industry"`

	SubscriptionPlan   string     `json:"subscription_plan"`
	SubscriptionStatus string     `json:"subscription_status"`
	SubscriptionEnd    *time.Time `json:"subscription_end"`
	FiturAktif         string     `json:"fitur_aktif"`
	QuotaTerminal      int        `json:"quota_terminal"`

	Alamat    string `json:"alamat"`
	Provinsi  string `json:"provinsi"`
	Kota      string `json:"kota"`
	Kecamatan string `json:"kecamatan"`
	Kelurahan string `json:"kelurahan"`
	KodePos   string `json:"kode_pos"`

	LogoURL       string  `json:"logo_url"`
	IsTaxActive   bool    `json:"is_tax_active"`
	PajakPersen   float64 `json:"pajak_persen"`
	ReceiptFooter string  `json:"receipt_footer"`

	PaymentType string `json:"payment_type"`
	QrisImage   string `json:"qris_image"`
	QrisName    string `json:"qris_name"`

	PrinterWidth string `json:"printer_width"`
	PrinterType  string `json:"printer_type"`

	WaToken           string `json:"wa_token"`
	MidtransServerKey string `json:"midtrans_server_key"`
	MidtransClientKey string `json:"midtrans_client_key"`
}

// ==========================================
// PENGATURAN TOKO (STORE SETTINGS)
// ==========================================

func (h *RetailHandler) GetStoreSettings(c *gin.Context) {
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
	if err := src.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan!"})
		return
	}

	// Masking data rahasia sebelum dilempar ke client
	resp := StoreSettingsResponse{
		ID:           store.ID,
		PublicID:     store.PublicID,
		NamaToko:     store.NamaToko,
		Telepon:      store.Telepon,
		BusinessType: store.BusinessType,
		Industry:     store.Industry,

		SubscriptionPlan:   store.SubscriptionPlan,
		SubscriptionStatus: store.SubscriptionStatus,
		SubscriptionEnd:    store.SubscriptionEnd,
		FiturAktif:         store.FiturAktif,
		QuotaTerminal:      store.QuotaTerminal,

		Alamat:    store.Alamat,
		Provinsi:  store.Provinsi,
		Kota:      store.Kota,
		Kecamatan: store.Kecamatan,
		Kelurahan: store.Kelurahan,
		KodePos:   store.KodePos,

		LogoURL:       store.LogoURL,
		IsTaxActive:   store.IsTaxActive,
		PajakPersen:   store.PajakPersen,
		ReceiptFooter: store.ReceiptFooter,

		PaymentType: store.PaymentType,
		QrisImage:   store.QrisImage,
		QrisName:    store.QrisName,

		PrinterWidth: store.PrinterWidth,
		PrinterType:  store.PrinterType,
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

func (h *RetailHandler) UpdateStoreSettings(c *gin.Context) {
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

	var store models.Store
	if err := src.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan!"})
		return
	}

	// Menggunakan c.GetPostForm() agar owner bisa mengosongkan field jika diinginkan
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

	// Validasi input token khusus: jangan timpa jika frontend mengirimkan kode samaran masking
	if v, ok := c.GetPostForm("wa_token"); ok {
		if v != "HAS_TOKEN_HIDDEN_BY_SYSTEM" {
			store.WaToken = strings.TrimSpace(v)
		}
	}
	if v, ok := c.GetPostForm("midtrans_server_key"); ok {
		if v != "HAS_KEY_HIDDEN_BY_SYSTEM" {
			store.MidtransServerKey = strings.TrimSpace(v)
		}
	}
	if v, ok := c.GetPostForm("midtrans_client_key"); ok {
		if v != "HAS_KEY_HIDDEN_BY_SYSTEM" {
			store.MidtransClientKey = strings.TrimSpace(v)
		}
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
		allowedMime := map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true}
		contentType := file.Header.Get("Content-Type")
		if !allowedMime[contentType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format logo harus JPG, PNG atau WEBP"})
			return
		}
		fileSrc, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file logo"})
			return
		}
		defer fileSrc.Close()
		remotePath := fmt.Sprintf("stores/%s/logo", store.PublicID)
		urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, contentType, bucketName, remotePath)
		if errUpload != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal upload logo"})
			return
		}
		store.LogoURL = urlResult
	}

	if c.PostForm("delete_qris") == "true" {
		store.QrisImage = ""
	} else if file, err := c.FormFile("qris"); err == nil {
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran QRIS maksimal 5 MB"})
			return
		}
		allowedMime := map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true}
		contentType := file.Header.Get("Content-Type")
		if !allowedMime[contentType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format QRIS harus JPG, PNG atau WEBP"})
			return
		}
		fileSrc, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file QRIS"})
			return
		}
		defer fileSrc.Close()
		remotePath := fmt.Sprintf("stores/%s/qris", store.PublicID)
		urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, contentType, bucketName, remotePath)
		if errUpload != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal upload QRIS"})
			return
		}
		store.QrisImage = urlResult
	}

	if err := src.DB.Save(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengaturan toko"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pengaturan toko berhasil diperbarui!", "data": store})
}

// SETUP WHATSAPP TOKEN UNTUK NOTIF PEMBAYARAN OWNER
func (h *RetailHandler) TestWhatsApp(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	var store models.Store
	if err := src.DB.First(&store, storeIDRaw).Error; err != nil {
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

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	c.JSON(resp.StatusCode, result)
}

// ==========================================
// SUBSCRIPTION & MIDTRANS (PAYMENTS)
// ==========================================

type UpgradeInput struct {
	PlanName string `json:"plan_name"`
}

func (h *RetailHandler) CreateUpgradePayment(c *gin.Context) {
	storeIDRaw, exists := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) {
	case float64:
		storeID = uint(v)
	case uint:
		storeID = v
	case int:
		storeID = uint(v)
	}

	if !exists || storeID == 0 {
		userIDRaw, _ := c.Get("user_id")
		var user models.User
		if err := src.DB.First(&user, userIDRaw).Error; err == nil && user.StoreID != nil {
			storeID = *user.StoreID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Infrastruktur ID toko belum siap dikonfigurasi ."})
			return
		}
	}

	var input UpgradeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter pilihan paket tidak valid"})
		return
	}

	var finalPrice int64 = 0
	targetPlan := strings.ToLower(strings.TrimSpace(input.PlanName))

	switch targetPlan {
	case "basic":
		finalPrice = 49000
	case "pro":
		finalPrice = 149000
	case "premium":
		finalPrice = 299000
	case "terminal tambahan":
		finalPrice = 50000
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tingkatan tier paket SaaS tidak terekam di sistem."})
		return
	}

	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox
	if os.Getenv("APP_ENV") == "production" {
		midtrans.Environment = midtrans.Production
	}

	planCode := strings.ReplaceAll(strings.ToUpper(targetPlan), " ", "")
	orderID := fmt.Sprintf("UPGRADE-TOKO-%d-%s-%d", storeID, planCode, time.Now().Unix())

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{OrderID: orderID, GrossAmt: finalPrice},
		Items:              &[]midtrans.ItemDetails{{ID: "SUB-" + planCode, Price: finalPrice, Qty: 1, Name: "Langganan " + input.PlanName}},
	}

	snapResp, err := snap.CreateTransaction(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal berkoordinasi dengan Payment Gateway"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": snapResp.Token, "order_id": orderID})
}

func (h *RetailHandler) MidtransWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid"})
		return
	}

	orderID, _ := payload["order_id"].(string)
	statusCode, _ := payload["status_code"].(string)
	grossAmount, _ := payload["gross_amount"].(string)
	signatureKey, _ := payload["signature_key"].(string)
	transactionStatus, _ := payload["transaction_status"].(string)

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	rawData := orderID + statusCode + grossAmount + serverKey
	sha := sha512.New()
	sha.Write([]byte(rawData))
	expectedSignature := hex.EncodeToString(sha.Sum(nil))

	if signatureKey != expectedSignature {
		fmt.Println("ALERT FRAUD: Signature Webhook Palsu/Mencurigakan!")
		c.JSON(http.StatusForbidden, gin.H{"error": "Banned! Data Signature tidak cocok!"})
		return
	}

	if transactionStatus == "settlement" || transactionStatus == "capture" {
		parts := strings.Split(orderID, "-")
		db := h.Repo.GetDB()

		if len(parts) >= 4 && parts[0] == "UPGRADE" {
			storeID := parts[2]
			planName := parts[3]

			if planName == "TERMINALTAMBAHAN" {
				db.Exec("UPDATE stores SET quota_terminal = quota_terminal + 1 WHERE id = ?", storeID)
			} else {
				endDate := time.Now().AddDate(0, 1, 0)
				db.Exec("UPDATE stores SET subscription_status = ?, subscription_end = ?, subscription_plan = ? WHERE id = ?", "active", endDate, strings.ToLower(planName), storeID)
			}
		} else if len(parts) >= 2 && parts[0] == "POS" {
			// SECURITY FIXED: Menggunakan safe parameter binding (?) untuk mencegah SQL Injection via orderID
			db.Exec("UPDATE transactions SET status_bayar = ? WHERE no_invoice = ?", "LUNAS", orderID)
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ==========================================
// MIDTRANS TRANSAKSI KASIR (UANG MASUK KE REKENING TENANT)
// ==========================================

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

	var s snap.Client
	env := midtrans.Sandbox
	if os.Getenv("APP_ENV") == "production" {
		env = midtrans.Production
	}

	s.New(store.MidtransServerKey, env)
	orderID := fmt.Sprintf("POS-STR%d-%d", storeID, time.Now().Unix())

	req := &snap.Request{TransactionDetails: midtrans.TransactionDetails{OrderID: orderID, GrossAmt: int64(input.Total)}}
	snapResp, err := s.CreateTransaction(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.GetMessage()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": snapResp.Token, "order_id": orderID})
}
