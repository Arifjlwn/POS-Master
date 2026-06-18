package delivery

import (
	"fmt"
	"net/http"
	"time"

	"pos-backend/models"
	"pos-backend/src/modules/jasalayanan/laundry/domain"
	"pos-backend/src/modules/jasalayanan/laundry/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 🚀 BIKIN HANDLER KHUSUS RAK (Sesuai Clean Architecture lu bray!)
type LaundryRackHandler struct {
	Repo repository.LaundryRepository
}

func NewLaundryRackHandler(repo repository.LaundryRepository) *LaundryRackHandler {
	return &LaundryRackHandler{Repo: repo}
}

type RackResponse struct {
	domain.LaundryRack
	IsiCucian    int                               `json:"isi_cucian"`
	DetailCucian []domain.TransactionLaundryDetail `json:"detail_cucian"`
}

// 1. 🎛️ GET ALL RACKS
func (h *LaundryRackHandler) GetRacks(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := extractUintID(storeIDRaw)

	db := h.Repo.GetDB()

	// Tarik fisik rak
	var racks []domain.LaundryRack
	if err := db.Where("store_id = ?", storeID).Order("baris ASC, kolom ASC").Find(&racks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengumpulkan data konfigurasi rak laundry"})
		return
	}

	// Tarik semua detail cucian yang belum diambil/selesai dan punya rack_id
	var activeDetails []domain.TransactionLaundryDetail
	db.Preload("Product").Where("store_id = ? AND status_cucian NOT IN ('DIAMBIL', 'SELESAI') AND rack_id IS NOT NULL", storeID).Find(&activeDetails)

	// Map data cucian ke rak masing-masing
	var result []RackResponse
	for _, r := range racks {
		rackRes := RackResponse{
			LaundryRack:  r,
			IsiCucian:    0,
			DetailCucian: []domain.TransactionLaundryDetail{},
		}

		for _, d := range activeDetails {
			if d.RackID == r.ID {
				rackRes.IsiCucian++
				rackRes.DetailCucian = append(rackRes.DetailCucian, d)
			}
		}
		result = append(result, rackRes)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   result,
	})
}

// 2. 🚀 SETUP INITIAL RACKS
type SetupRackInput struct {
	JumlahBaris int `json:"jumlah_baris" binding:"required,gt=0,lte=26"`
	JumlahKolom int `json:"jumlah_kolom" binding:"required,gt=0,lte=50"`
}

func (h *LaundryRackHandler) SetupInitialRacks(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := extractUintID(storeIDRaw)

	var input SetupRackInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input jumlah baris atau kolom tidak valid bray!"})
		return
	}

	db := h.Repo.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		var count int64
		tx.Model(&domain.LaundryRack{}).Where("store_id = ?", storeID).Count(&count)
		if count > 0 {
			return fmt.Errorf("ruko Anda sudah memiliki konfigurasi rak aktif, gunakan menu upgrade kapasitas")
		}

		var racksToInsert []domain.LaundryRack
		for b := 1; b <= input.JumlahBaris; b++ {
			hurufBaris := string(rune(64 + b))
			for k := 1; k <= input.JumlahKolom; k++ {
				namaRak := fmt.Sprintf("%s-%d", hurufBaris, k)
				racksToInsert = append(racksToInsert, domain.LaundryRack{
					StoreID:   storeID,
					NamaRak:   namaRak,
					Baris:     b,
					Kolom:     k,
					Status:    "TERSEDIA",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			}
		}

		if len(racksToInsert) > 0 {
			if err := tx.Create(&racksToInsert).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "BOOM! Konfigurasi susunan rak baju berhasil di-generate otomatis!"})
}

// 3. 🛠️ TOGGLE RACK STATUS
type UpdateRackStatusInput struct {
	Status string `json:"status" binding:"required"`
}

func (h *LaundryRackHandler) ToggleRackStatus(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := extractUintID(storeIDRaw)
	rackID := c.Param("id")

	var input UpdateRackStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status wajib diisi (TERSEDIA / RUSAK)"})
		return
	}

	db := h.Repo.GetDB()
	var rack domain.LaundryRack

	if err := db.Where("id = ? AND store_id = ?", rackID, storeID).First(&rack).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target fisik rak tidak ditemukan bray"})
		return
	}

	rack.Status = input.Status
	rack.UpdatedAt = time.Now()

	if err := db.Save(&rack).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-update status operasional rak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": fmt.Sprintf("Rak %s berhasil diset menjadi %s", rack.NamaRak, rack.Status)})
}

// 4. 🔄 CHANGE ORDER RACK
type PindahRakInput struct {
	NewRackID uint   `json:"new_rack_id" binding:"required"`
	Invoice   string `json:"invoice" binding:"required"`
}

func (h *LaundryRackHandler) ChangeOrderRack(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := extractUintID(storeIDRaw)

	var input PindahRakInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target rak dan nomor invoice wajib dilampirkan"})
		return
	}

	db := h.Repo.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Validasi tujuan rak baru apakah eksis and sehat walafiat
		var newRack domain.LaundryRack
		if err := tx.Where("id = ? AND store_id = ? AND status = 'TERSEDIA'", input.NewRackID, storeID).First(&newRack).Error; err != nil {
			return fmt.Errorf("target rak baru tidak aktif atau sedang rusak")
		}

		// 2. Cari transaksi induk berdasarkan No Invoice
		var trx models.Transaction
		if err := tx.Where("no_invoice = ? AND store_id = ?", input.Invoice, storeID).First(&trx).Error; err != nil {
			return fmt.Errorf("berkas transaksi dengan invoice %s tidak ditemukan", input.Invoice)
		}

		// 3. 🚀 UPDATE MASSAL: Pindahin SEMUA cucian yang ada di invoice ini ke rak baru!
		if err := tx.Model(&domain.TransactionLaundryDetail{}).
			Where("transaction_id = ? AND store_id = ?", trx.ID, storeID).
			Updates(map[string]interface{}{
				"rack_id":    newRack.ID,
				"nomor_rak":  newRack.NamaRak,
				"updated_at": time.Now(),
			}).Error; err != nil {
			return fmt.Errorf("gagal memindahkan item cucian ke rak baru")
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "Seluruh cucian pada nota ini berhasil dievakuasi ke rak baru!"})
}

// --- HELPER UNTUK EXTRAKSI TOKEN STORE ID AMAN BRAY ---
func extractUintID(raw interface{}) uint {
	switch v := raw.(type) {
	case float64:
		return uint(v)
	case uint:
		return v
	case int:
		return uint(v)
	default:
		return 0
	}
}
