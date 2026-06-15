package delivery

import (
	"net/http"
	"strconv"
	"pos-backend/src/modules/jasalayanan/laundry/domain"
	"pos-backend/src/modules/jasalayanan/laundry/repository"
	"github.com/gin-gonic/gin"
)

type LaundryPerfumeHandler struct {
	Repo repository.LaundryRepository
}

func NewLaundryPerfumeHandler(r repository.LaundryRepository) *LaundryPerfumeHandler {
	return &LaundryPerfumeHandler{Repo: r}
}

type PerfumeInput struct {
	Nama  string  `json:"nama" binding:"required"`
	Harga float64 `json:"harga"`
}

func (h *LaundryPerfumeHandler) GetPerfumes(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	perfumes, err := h.Repo.GetPerfumesByStoreID(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data parfum"})
		return
	}
	c.JSON(http.StatusOK, perfumes)
}

func (h *LaundryPerfumeHandler) CreatePerfume(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	var input PerfumeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid"})
		return
	}
	newPerfume := domain.Perfume{StoreID: storeID, Nama: input.Nama, Harga: input.Harga, Status: "AKTIF"}
	if err := h.Repo.CreatePerfume(&newPerfume); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan varian parfum baru"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Parfum berhasil ditambahkan!", "data": newPerfume})
}

func (h *LaundryPerfumeHandler) DeletePerfume(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	perfumeID, _ := strconv.Atoi(c.Param("id"))
	if err := h.Repo.DeletePerfume(uint(perfumeID), storeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus parfum"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Varian parfum berhasil dihapus!"})
}