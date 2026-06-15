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
	NamaProduk  string  `json:"nama_produk" binding:"required"`
	HargaJual   float64 `json:"harga_jual" binding:"required"`
	SatuanDasar string  `json:"satuan_dasar" binding:"required"` 
	Estimasi    string  `json:"estimasi"`                                       
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
	newLayanan := models.Product{StoreID: storeID, NamaProduk: input.NamaProduk, Kategori: "JASA_LAUNDRY", HargaJual: input.HargaJual, SatuanDasar: input.SatuanDasar, Estimasi: input.Estimasi, Stok: 0}
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
	layanan.Estimasi = input.Estimasi
	if err := h.Repo.UpdateLayanan(layanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui layanan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil diperbarui!"})
}

func (h *LaundryServiceHandler) HapusLayananLaundry(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	productID, _ := strconv.Atoi(c.Param("id"))
	if err := h.Repo.DeleteLayanan(uint(productID), storeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus layanan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil dihapus"})
}