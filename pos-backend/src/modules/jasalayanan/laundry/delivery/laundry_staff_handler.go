package delivery

import (
	"net/http"
	"strconv"
	"strings"
	"pos-backend/models"
	"pos-backend/src/modules/jasalayanan/laundry/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LaundryStaffHandler struct {
	Repo repository.LaundryRepository
}

func NewLaundryStaffHandler(r repository.LaundryRepository) *LaundryStaffHandler {
	return &LaundryStaffHandler{Repo: r}
}

type KasirInput struct {
	Name     string `json:"name" binding:"required"`
	NoHP     string `json:"no_hp"`
	Password string `json:"password"`
}

func (h *LaundryStaffHandler) GetKasirList(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	kasirList, err := h.Repo.GetKasirByStoreID(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data kasir"})
		return
	}
	var response []map[string]interface{}
	for _, k := range kasirList {
		email := ""
		if k.Email != nil { email = *k.Email }
		response = append(response, map[string]interface{}{"id": k.ID, "name": k.Name, "email": email, "no_hp": k.NoHP})
	}
	c.JSON(http.StatusOK, response)
}

func (h *LaundryStaffHandler) CreateKasir(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	var input KasirInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak lengkap"})
		return
	}
	store, err := h.Repo.GetStoreByID(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonfirmasi infrastruktur toko"})
		return
	}
	namaKasirBersih := strings.ToLower(strings.ReplaceAll(input.Name, " ", ""))
	namaTokoBersih := strings.ToLower(strings.ReplaceAll(store.NamaToko, " ", ""))
	emailDummy := "kasir." + namaKasirBersih + "@" + namaTokoBersih + ".com"
	passToHash := "kasir123"
	if input.Password != "" { passToHash = input.Password }
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(passToHash), bcrypt.DefaultCost)
	newKasir := models.User{Name: input.Name, Email: &emailDummy, Password: string(hashedPassword), Role: "kasir", StoreID: &storeID, NoHP: input.NoHP, IsVerified: true}
	if err := h.Repo.CreateKasir(&newKasir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan kasir. Pastikan nama tidak duplikat."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Kasir berhasil ditambahkan!", "email": emailDummy})
}

func (h *LaundryStaffHandler) DeleteKasir(c *gin.Context) {
	storeIDRaw, _ := c.Get("store_id")
	storeID := uint(storeIDRaw.(float64))
	kasirID, _ := strconv.Atoi(c.Param("id"))
	if err := h.Repo.DeleteKasir(uint(kasirID), storeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data kasir"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Karyawan berhasil dihapus dari sistem!"})
}