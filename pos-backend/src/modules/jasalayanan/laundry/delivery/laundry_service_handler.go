package delivery

import (
	"net/http"
	"strconv"
	"pos-backend/models"
	"pos-backend/src/modules/jasalayanan/laundry/repository"
	"github.com/gin-gonic/gin"
)

type LaundryServiceHandler struct {
	Repo repository.LaundryRepository
}

func NewLaundryServiceHandler(r repository.LaundryRepository) *LaundryServiceHandler {
	return &LaundryServiceHandler{Repo: r}
}

type InputLayanan struct {
	NamaProduk     string  `json:"nama_produk" binding:"required"`
	SatuanDasar    string  `json:"satuan_dasar" binding:"required"`
	EstimasiDurasi int     `json:"estimasi_durasi" binding:"required,gt=0"`
	EstimasiSatuan string  `json:"estimasi_satuan" binding:"required"`
	HargaJual      float64 `json:"harga_jual" binding:"required"`                                       
}

func (h *LaundryServiceHandler) AmbilDaftarLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	listJasa, err := h.Repo.GetLayananLaundry(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil katalog layanan laundry"})
		return
	}
	c.JSON(http.StatusOK, listJasa)
}

func (h *LaundryServiceHandler) TambahLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	var input InputLayanan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input layanan tidak valid"})
		return
	}
	newLayanan := models.Product{ // atau nama model lu
		StoreID:        storeID,
		NamaProduk:     input.NamaProduk,
		SatuanDasar:    input.SatuanDasar,
		HargaJual:      input.HargaJual,
		EstimasiDurasi: input.EstimasiDurasi,
		EstimasiSatuan: input.EstimasiSatuan,
		ProductType:     "JASA_LAUNDRY",
	}
	if err := h.Repo.CreateLayanan(&newLayanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan layanan baru"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil ditambahkan", "data": newLayanan})
}

func (h *LaundryServiceHandler) EditLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	productID, _ := strconv.Atoi(c.Param("id"))
	var input InputLayanan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input layanan tidak valid"})
		return
	}
	layanan, err := h.Repo.GetLayananByID(uint(productID), storeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Layanan tidak ditemukan"})
		return
	}
	layanan.NamaProduk = input.NamaProduk
	layanan.HargaJual = input.HargaJual
	layanan.SatuanDasar = input.SatuanDasar
	layanan.EstimasiDurasi = input.EstimasiDurasi
	layanan.EstimasiSatuan = input.EstimasiSatuan

	if err := h.Repo.UpdateLayanan(layanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui layanan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil diperbarui!"})
}

func (h *LaundryServiceHandler) HapusLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64)) // ⚠️ Gunakan helper lu kalau ada, atau biarkan pakai ini
	
	productID, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.Repo.DeleteLayanan(uint(productID), storeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 🚀 Ubah pesannya biar elegan
	c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "Layanan berhasil dihapus / diarsipkan dari katalog"})
}