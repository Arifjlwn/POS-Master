package delivery

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

// ==========================================
// 💳 SUBSCRIPTION & MIDTRANS (PAYMENTS)
// ==========================================
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