package delivery

import (
	"pos-backend/src/core/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRetailInventoryRoutes(rg *gin.RouterGroup, h *RetailHandler) {

	// =====================================
	// 🟢 FITUR GERBANG UTAMA (BYPASS BILLING CHECK FOR PAYMENT GATEWAY)
	// =====================================

	// 🚀 FIX MUTLAK: Cabut RequireSaaSLevel dari rute upgrade pembayaran bray! Toko inactive WAJIB bisa ngakses rute ini buat dapet Snap Token Midtrans!
	rg.POST("/subscription/upgrade", h.CreateUpgradePayment)

	// =====================================
	// 🟢 FITUR LEVEL 1 (BASIC & CORE OPERATIONAL RETAIL)
	// =====================================

	rg.POST("/pos/midtrans-order", middlewares.RequireSaaSLevel(1), h.CreatePosMidtransOrder)

	// Master Produk & Kategori Catalog
	rg.POST("/products", middlewares.RequireSaaSLevel(1), h.CreateProduct)
	rg.GET("/products", middlewares.RequireSaaSLevel(1), h.GetProducts)
	rg.PUT("/products/:id", middlewares.RequireSaaSLevel(1), h.UpdateProduct)
	rg.DELETE("/products/:id", middlewares.RequireSaaSLevel(1), h.DeleteProduct)
	rg.GET("/categories", middlewares.RequireSaaSLevel(1), h.GetCategories)
	rg.GET("/products/export", middlewares.RequireSaaSLevel(1), h.ExportProducts)
	rg.POST("/products/import", middlewares.RequireSaaSLevel(1), h.ImportProducts)

	// Laporan Penerimaan Barang / LPB (Supplier)
	rg.POST("/purchases", middlewares.RequireSaaSLevel(1), h.CreateLPB)

	// Shift POS / Cashier Session
	rg.POST("/pos/open-session", middlewares.RequireSaaSLevel(1), h.OpenSession)
	rg.GET("/pos/check-session", middlewares.RequireSaaSLevel(1), h.CheckSessionStatus)
	rg.POST("/pos/close-session/:id", middlewares.RequireSaaSLevel(1), h.CloseSession)

	// POS Checkout Transaksi Kasir
	rg.POST("/checkout", middlewares.RequireSaaSLevel(1), h.CreateTransaction)
	rg.GET("/transactions", middlewares.RequireSaaSLevel(1), h.GetTransactions)
	rg.GET("/journal/closing", middlewares.RequireSaaSLevel(1), h.GetDailyClosing)

	// Store Profiles Settings
	rg.GET("/store/settings", middlewares.RequireSaaSLevel(1), h.GetStoreSettings)
	rg.PUT("/store/settings", middlewares.RequireSaaSLevel(1), h.UpdateStoreSettings)

	// =====================================
	// 🟡 FITUR LEVEL 2 (PRO - TIM & ABSENSI)
	// =====================================

	// Absensi Karyawan
	rg.POST("/attendance", middlewares.RequireSaaSLevel(2), h.StoreAttendance)
	rg.GET("/attendance", middlewares.RequireSaaSLevel(2), h.GetAttendance)
	rg.GET("/attendance/export", middlewares.RequireSaaSLevel(2), h.ExportAttendance)

	// HR / Data Karyawan
	rg.POST("/employees", middlewares.RequireSaaSLevel(2), h.CreateEmployee)
	rg.GET("/employees", middlewares.RequireSaaSLevel(2), h.GetEmployees)
	rg.PUT("/employees/:id", middlewares.RequireSaaSLevel(2), h.UpdateEmployee)
	rg.DELETE("/employees/:id", middlewares.RequireSaaSLevel(2), h.DeleteEmployee)

	// Jadwal Kerja / Rostering
	rg.POST("/schedules/bulk", middlewares.RequireSaaSLevel(2), h.SaveSchedules)
	rg.GET("/schedules", middlewares.RequireSaaSLevel(2), h.GetSchedules)

	// =====================================
	// 🔴 FITUR LEVEL 3 (PREMIUM - ANALYTICS & INVENTORY AUDIT)
	// =====================================

	// Dashboard Analytics Report Owner
	rg.GET("/report/dashboard", middlewares.RequireSaaSLevel(3), h.GetDashboardReport)

	// Stock Opname
	rg.POST("/stock-opname", middlewares.RequireSaaSLevel(3), h.CreateStockOpname)
	rg.GET("/stock-opname/last-status", middlewares.RequireSaaSLevel(3), h.GetLastSOStatus)
	rg.GET("/stock-opname/last-minus", middlewares.RequireSaaSLevel(3), h.GetLastSOMinusItems)
	rg.GET("/stock-opname/history", middlewares.RequireSaaSLevel(3), h.GetStockOpnameHistory)
	
	// 🛡️ FIX MULTIPLE CRITICAL: Ubah PATCH -> POST & sesuaikan path urutan agar klop dengan Axios Vue bray!
	rg.POST("/stock-opname/approve/:id", middlewares.RequireSaaSLevel(3), h.ApproveStockOpname)

	// Klaim Penemuan Barang Nyempil (Adjustment)
	rg.POST("/stock-adjustment/request", middlewares.RequireSaaSLevel(3), h.SubmitKlaimBarang)
	rg.GET("/stock-adjustment/history", middlewares.RequireSaaSLevel(3), h.GetStockAdjustmentHistory)
	
	// 🛡️ FIX MULTIPLE CRITICAL: Ubah PATCH -> POST & sesuaikan path urutan approval klaim bray bray!
	rg.POST("/stock-adjustment/approve/:id", middlewares.RequireSaaSLevel(3), h.ApproveStockAdjustment)

	// Retur Barang
	rg.POST("/returns", middlewares.RequireSaaSLevel(3), h.CreateReturn)
	rg.GET("/returns", middlewares.RequireSaaSLevel(3), h.GetReturns)

	// INTEGRASI EKSTERNAL
	rg.POST("/whatsapp/test", middlewares.RequireSaaSLevel(3), h.TestWhatsApp)
}
