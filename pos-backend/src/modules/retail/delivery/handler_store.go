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

	"pos-backend/models"
	src "pos-backend/src/core/config"
	"pos-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// ==========================================
// ⚙️ PENGATURAN TOKO (STORE SETTINGS)
// ==========================================

func (h *RetailHandler) GetStoreSettings(c *gin.Context) {
	storeIDRaw, exists := c.Get("store_id")
	if !exists { c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"}); return }
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	var store models.Store
	if err := src.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": store})
}

func (h *RetailHandler) UpdateStoreSettings(c *gin.Context) {
	roleOwner, _ := c.Get("role")
	if roleOwner != "owner" && roleOwner != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Owner atau Manager yang bisa mengubah pengaturan toko!"})
		return
	}

	storeIDRaw, _ := c.Get("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	var store models.Store
	if err := src.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan!"})
		return
	}

	if v := c.PostForm("nama_toko"); v != "" { store.NamaToko = v }
	if v := c.PostForm("telepon"); v != "" { store.Telepon = utils.FormatPhoneNumber(v) }
	if v := c.PostForm("alamat"); v != "" { store.Alamat = v }
	if v := c.PostForm("provinsi"); v != "" { store.Provinsi = v }
	if v := c.PostForm("kota"); v != "" { store.Kota = v }
	if v := c.PostForm("kecamatan"); v != "" { store.Kecamatan = v }
	if v := c.PostForm("kelurahan"); v != "" { store.Kelurahan = v }
	if v := c.PostForm("kode_pos"); v != "" { store.KodePos = v }
	if v := c.PostForm("qris_name"); v != "" { store.QrisName = v }
	if v := c.PostForm("receipt_footer"); v != "" { store.ReceiptFooter = v }
	if v := c.PostForm("wa_token"); v != "" { store.WaToken = v }
	if v := c.PostForm("payment_type"); v != "" { store.PaymentType = v }
	if v := c.PostForm("midtrans_server_key"); v != "" { store.MidtransServerKey = v }
	if v := c.PostForm("midtrans_client_key"); v != "" { store.MidtransClientKey = v }
	if v := c.PostForm("printer_width"); v != "" { store.PrinterWidth = v }
	if v := c.PostForm("printer_type"); v != "" { store.PrinterType = v }

	if v := c.PostForm("is_tax_active"); v != "" { store.IsTaxActive = (v == "true") }
	if v := c.PostForm("pajak_persen"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil { store.PajakPersen = parsed }
	}

	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")

	if c.PostForm("delete_logo") == "true" {
		store.LogoURL = ""
	} else if file, err := c.FormFile("logo"); err == nil {
		if file.Size > 5*1024*1024 { c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran logo maksimal 5 MB"}); return }
		allowedMime := map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true}
		contentType := file.Header.Get("Content-Type")
		if !allowedMime[contentType] { c.JSON(http.StatusBadRequest, gin.H{"error": "Format logo harus JPG, PNG atau WEBP"}); return }
		fileSrc, err := file.Open(); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file logo"}); return }
		defer fileSrc.Close()
		remotePath := fmt.Sprintf("stores/%d/logo", storeID)
		urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, contentType, bucketName, remotePath)
		if errUpload != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal upload logo"}); return }
		store.LogoURL = urlResult
	}

	if c.PostForm("delete_qris") == "true" {
		store.QrisImage = ""
	} else if file, err := c.FormFile("qris"); err == nil {
		if file.Size > 5*1024*1024 { c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran QRIS maksimal 5 MB"}); return }
		allowedMime := map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true}
		contentType := file.Header.Get("Content-Type")
		if !allowedMime[contentType] { c.JSON(http.StatusBadRequest, gin.H{"error": "Format QRIS harus JPG, PNG atau WEBP"}); return }
		fileSrc, err := file.Open(); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca file QRIS"}); return }
		defer fileSrc.Close()
		remotePath := fmt.Sprintf("stores/%d/qris", storeID)
		urlResult, errUpload := utils.UploadToSupabase(fileSrc, file.Filename, contentType, bucketName, remotePath)
		if errUpload != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal upload QRIS"}); return }
		store.QrisImage = urlResult
	}

	if err := src.DB.Save(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengaturan toko"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pengaturan toko berhasil diperbarui!", "data": store})
}

// ==========================================
// 💳 SUBSCRIPTION & MIDTRANS (PAYMENTS)
// ==========================================

type UpgradeInput struct {
	PlanName string `json:"plan_name"`
	Price    int64  `json:"price"`
}

func (h *RetailHandler) CreateUpgradePayment(c *gin.Context) {
	storeIDRaw := c.MustGet("store_id")
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	var input UpgradeInput
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Data paket tidak valid"}); return }

	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox
	if os.Getenv("APP_ENV") == "production" { midtrans.Environment = midtrans.Production }

	planCode := strings.ReplaceAll(strings.ToUpper(input.PlanName), " ", "")
	orderID := fmt.Sprintf("UPGRADE-TOKO-%d-%s-%d", storeID, planCode, time.Now().Unix())

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{OrderID: orderID, GrossAmt: input.Price},
		Items: &[]midtrans.ItemDetails{{ID: "SUB-" + planCode, Price: input.Price, Qty: 1, Name: "Langganan Paket " + input.PlanName}},
	}

	snapResp, err := snap.CreateTransaction(req)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungi Payment Gateway"}); return }
	c.JSON(http.StatusOK, gin.H{"token": snapResp.Token, "order_id": orderID})
}

func (h *RetailHandler) MidtransWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid"}); return }

	orderID, _ := payload["order_id"].(string)
	statusCode, _ := payload["status_code"].(string)
	grossAmount, _ := payload["gross_amount"].(string)
	signatureKey, _ := payload["signature_key"].(string)
	transactionStatus, _ := payload["transaction_status"].(string)

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	rawData := orderID + statusCode + grossAmount + serverKey
	sha := sha512.New(); sha.Write([]byte(rawData))
	expectedSignature := hex.EncodeToString(sha.Sum(nil))

	if signatureKey != expectedSignature {
		fmt.Println("❌ ALERT FRAUD: Signature Webhook Palsu/Mencurigakan!")
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
			db.Exec("UPDATE transactions SET status_bayar = ? WHERE no_invoice = ?", "LUNAS", orderID)
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ==========================================
// 🚀 MIDTRANS TRANSAKSI KASIR (UANG MASUK KE REKENING TENANT)
// ==========================================

type PosMidtransReq struct { Total float64 `json:"total" binding:"required"` }

func (h *RetailHandler) CreatePosMidtransOrder(c *gin.Context) {
	storeIDRaw, exists := c.Get("store_id")
	if !exists { c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"}); return }
	var storeID uint
	switch v := storeIDRaw.(type) { case float64: storeID = uint(v); case uint: storeID = v; case int: storeID = uint(v) }

	var input PosMidtransReq
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Data nominal tidak valid"}); return }

	var store models.Store
	if err := src.DB.First(&store, storeID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan!"}); return }
	if store.PaymentType != "midtrans" || store.MidtransServerKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Toko belum mengatur Midtrans Server Key!"})
		return
	}

	var s snap.Client
	env := midtrans.Sandbox
	if os.Getenv("APP_ENV") == "production" { env = midtrans.Production }

	s.New(store.MidtransServerKey, env)
	orderID := fmt.Sprintf("POS-STR%d-%d", storeID, time.Now().Unix())

	req := &snap.Request{TransactionDetails: midtrans.TransactionDetails{OrderID: orderID, GrossAmt: int64(input.Total)}}
	snapResp, err := s.CreateTransaction(req)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.GetMessage()}); return }

	c.JSON(http.StatusOK, gin.H{"token": snapResp.Token, "order_id": orderID})
}