package delivery

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"pos-backend/models"
	"pos-backend/src/core/repository"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type BillingHandler struct {
	Repo repository.CoreRepo
}

func NewBillingHandler(repo repository.CoreRepo) *BillingHandler {
	return &BillingHandler{Repo: repo}
}

type UpgradeInput struct {
	PlanName string `json:"plan_name"`
}

func (h *BillingHandler) CreateUpgradePayment(c *gin.Context) {
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

	db := h.Repo.GetDB()
	if !exists || storeID == 0 {
		userIDRaw, _ := c.Get("user_id")
		var user models.User
		if err := db.First(&user, userIDRaw).Error; err == nil && user.StoreID != nil {
			storeID = *user.StoreID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Infrastruktur ID toko belum siap dikonfigurasi."})
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

func (h *BillingHandler) MidtransWebhook(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Signature tidak cocok!"})
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
