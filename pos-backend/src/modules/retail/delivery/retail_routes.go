package delivery

import (
	"pos-backend/src/core/middlewares" // 🚀 IMPORT MIDDLEWARE LU DI SINI

	"github.com/gin-gonic/gin"
)

// Master Produk CRUD
func RegisterRetailInventoryRoutes(rg *gin.RouterGroup, h *RetailHandler) {

	// =====================================
	// 🟢 FITUR MIDTRANS USER
	// =====================================
	rg.POST("/pos/midtrans-order", h.CreatePosMidtransOrder)

	// =====================================
	// 🟢 FITUR LEVEL 1 (BASIC)
	// =====================================
	
	// Master Produk & Kategori
	rg.POST("/products", h.CreateProduct)
	rg.GET("/products", h.GetProducts)
	rg.PUT("/products/:id", h.UpdateProduct)
	rg.DELETE("/products/:id", h.DeleteProduct)
	rg.GET("/categories", h.GetCategories)
	
	// Impor & Ekspor Barang
	rg.GET("/products/export", h.ExportProducts)
	rg.POST("/products/import", h.ImportProducts)

	// Laporan Penerimaan Barang / LPB (Bisa dipake Basic)
	rg.POST("/purchases", h.CreateLPB)

	// Shift POS / Cashier Session (Open-Close Kasir)
	rg.POST("/pos/open-session", h.OpenSession)
	rg.GET("/pos/check-session", h.CheckSessionStatus)
	rg.POST("/pos/close-session/:id", h.CloseSession)

	// POS Checkout Transaksi & Struk Nota
	rg.POST("/checkout", h.CreateTransaction)
	rg.GET("/transactions", h.GetTransactions)

	// Setting Toko & Akun
	rg.GET("/store/settings", h.GetStoreSettings)
	rg.PUT("/store/settings", h.UpdateStoreSettings)
	rg.POST("/subscription/upgrade", h.CreateUpgradePayment)


	// =====================================
	// 🟡 FITUR LEVEL 2 (PRO)
	// =====================================
	
	// Absensi Karyawan
	rg.POST("/attendance", middlewares.RequireSaaSLevel(2), h.StoreAttendance)
	rg.GET("/attendance", middlewares.RequireSaaSLevel(2), h.GetAttendance)
	rg.GET("/attendance/export", middlewares.RequireSaaSLevel(2), h.ExportAttendance)

	// HR / Data Karyawan
	rg.POST("/employees", middlewares.RequireSaaSLevel(2), h.CreateEmployee)
	rg.GET("/employees", middlewares.RequireSaaSLevel(2), h.GetEmployees)
	rg.PUT("/employees/:id", middlewares.RequireSaaSLevel(2), h.UpdateEmployee)

	// Jadwal Kerja / Rostering Karyawan
	rg.POST("/schedules/bulk", middlewares.RequireSaaSLevel(2), h.SaveSchedules)
	rg.GET("/schedules", middlewares.RequireSaaSLevel(2), h.GetSchedules)


	// =====================================
	// 🔴 FITUR LEVEL 3 (PREMIUM)
	// =====================================
	
	// Dashboard Analytics Report Owner
	rg.GET("/report/dashboard", middlewares.RequireSaaSLevel(3), h.GetDashboardReport)

	// Stock Opname
	rg.POST("/stock-opname", middlewares.RequireSaaSLevel(3), h.CreateStockOpname)
	rg.GET("/stock-opname/last-status", middlewares.RequireSaaSLevel(3), h.GetLastSOStatus)
	rg.GET("/stock-opname/last-minus", middlewares.RequireSaaSLevel(3), h.GetLastSOMinusItems)
	rg.GET("/stock-opname/history", middlewares.RequireSaaSLevel(3), h.GetStockOpnameHistory)
	rg.PATCH("/stock-opname/:id/approve", middlewares.RequireSaaSLevel(3), h.ApproveStockOpname)
	
	rg.POST("/stock-adjustment/request", middlewares.RequireSaaSLevel(3), h.SubmitKlaimBarang)
	rg.GET("/stock-adjustment/history", middlewares.RequireSaaSLevel(3), h.GetStockAdjustmentHistory)
	rg.PATCH("/stock-adjustment/:id/approve", middlewares.RequireSaaSLevel(3), h.ApproveStockAdjustment)

	// Retur Barang
	rg.POST("/returns", middlewares.RequireSaaSLevel(3), h.CreateReturn)
	rg.GET("/returns", middlewares.RequireSaaSLevel(3), h.GetReturns)
}