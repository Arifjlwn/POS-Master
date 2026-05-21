package delivery

import "github.com/gin-gonic/gin"

func RegisterLaundryRoutes(rg *gin.RouterGroup, h *LaundryHandler) {
	// Kasir Karyawan Management
	rg.GET("/kasir", h.GetKasirList)
	rg.POST("/kasir", h.CreateKasir)
	rg.DELETE("/kasir/:id", h.DeleteKasir)

	// Kasir POS Laundry & Live Search
	rg.GET("/services", h.AmbilDaftarLayananLaundry)
	rg.POST("/checkout", h.ProsesCheckoutLaundry)
	rg.GET("/customers/search", h.CariPelanggan)

	// Pelunasan Piutang / Bill
	rg.PUT("/transactions/:id/lunas", h.LunasiTransaksi)

	// Master Paket / Jasa Laundry
	rg.POST("/services/new", h.TambahLayananLaundry) // Diubah dikit pathnya biar gampang dibedain sama GET services
	rg.PUT("/services/:id", h.EditLayananLaundry)
	rg.DELETE("/services/:id", h.HapusLayananLaundry)

	// Perfume Management
	rg.GET("/perfumes", h.GetPerfumes)
	rg.POST("/perfumes", h.CreatePerfume)
	rg.DELETE("/perfumes/:id", h.DeletePerfume)

	// Report Keuangan & Antrean Cucian
	rg.GET("/report", h.GetLaporan)

	// TRACKING & SETTING TOKO
	rg.GET("/setting", h.GetSettingToko)
	rg.PUT("/setting", h.UpdateSettingToko)
	rg.GET("/tracking", h.AmbilDataTracking)
	rg.PUT("/tracking/:id/status", h.UpdateStatusCucian)
}