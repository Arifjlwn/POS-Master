package delivery

import (
	"pos-backend/src/modules/jasalayanan/laundry/repository"
	"pos-backend/src/modules/jasalayanan/laundry/usecase"

	"github.com/gin-gonic/gin"
)

func RegisterLaundryRoutes(rg *gin.RouterGroup, repo repository.LaundryRepository) {
	laundryUC := usecase.NewLaundryUseCase(repo)

	txHandler := NewLaundryTransactionHandler(laundryUC)
	serviceHandler := NewLaundryServiceHandler(repo)
	perfumeHandler := NewLaundryPerfumeHandler(repo)
	staffHandler := NewLaundryStaffHandler(repo)
	reportHandler := NewLaundryReportHandler(laundryUC, repo)

	// Staff
	rg.GET("/kasir", staffHandler.GetKasirList)
	rg.POST("/kasir", staffHandler.CreateKasir)
	rg.DELETE("/kasir/:id", staffHandler.DeleteKasir)

	// POS Core & Customer Search
	rg.GET("/services", serviceHandler.AmbilDaftarLayananLaundry)
	rg.POST("/checkout", txHandler.ProsesCheckoutLaundry)
	rg.GET("/customers/search", reportHandler.CariPelanggan)
	rg.PUT("/transactions/:id/lunas", txHandler.LunasiTransaksi)

	// CRUD Services
	rg.POST("/services", serviceHandler.TambahLayananLaundry)
	rg.PUT("/services/:id", serviceHandler.EditLayananLaundry)
	rg.DELETE("/services/:id", serviceHandler.HapusLayananLaundry)

	// Perfumes
	rg.GET("/perfumes", perfumeHandler.GetPerfumes)
	rg.POST("/perfumes", perfumeHandler.CreatePerfume)
	rg.DELETE("/perfumes/:id", perfumeHandler.DeletePerfume)

	// Reports & Settings
	rg.GET("/report", txHandler.GetLaporan)
	rg.GET("/tracking", reportHandler.AmbilDataTracking)
	rg.PUT("/tracking/:id/status", txHandler.UpdateStatusCucian)
}
