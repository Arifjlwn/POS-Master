package delivery

import "github.com/gin-gonic/gin"

func RegisterRetailInventoryRoutes(rg *gin.RouterGroup, h *RetailHandler) {
	// Master Produk CRUD
	rg.POST("/products", h.CreateProduct)
	rg.GET("/products", h.GetProducts)
	rg.PUT("/products/:id", h.UpdateProduct)
	rg.DELETE("/products/:id", h.DeleteProduct)
	rg.GET("/categories", h.GetCategories)
	
	// Impor & Ekspor Barang
	rg.GET("/products/export", h.ExportProducts)
	rg.POST("/products/import", h.ImportProducts)

	// Stock Opname
	rg.POST("/stock-opname", h.CreateStockOpname)
	rg.GET("/stock-opname/history", h.GetStockOpnameHistory)
	rg.PATCH("/stock-opname/:id/approve", h.ApproveStockOpname)

	// Retur Barang
	rg.POST("/returns", h.CreateReturn)
	rg.GET("/returns", h.GetReturns)

	// Laporan Penerimaan Barang / LPB
	rg.POST("/purchases", h.CreateLPB)

	// Absensi Karyawan
	rg.POST("/attendance", h.StoreAttendance)
	rg.GET("/attendance", h.GetAttendance)
	rg.GET("/attendance/export", h.ExportAttendance)

	// HR / Data Karyawan
	rg.POST("/employees", h.CreateEmployee)
	rg.GET("/employees", h.GetEmployees)
	rg.PUT("/employees/:id", h.UpdateEmployee)

	// Dashboard Analytics Report Owner
	rg.GET("/report/dashboard", h.GetDashboardReport)

	// 📅 Jadwal Kerja / Rostering Karyawan
	rg.POST("/schedules/bulk", h.SaveSchedules)
	rg.GET("/schedules", h.GetSchedules)

	// 💵 Shift POS / Cashier Session (Open-Close Kasir)
	rg.POST("/pos/open-session", h.OpenSession)
	rg.GET("/pos/check-session", h.CheckSessionStatus)
	rg.POST("/pos/close-session/:id", h.CloseSession)

	// 🛒 POS Checkout Transaksi & Struk Nota
	rg.POST("/checkout", h.CreateTransaction)
	rg.GET("/transactions", h.GetTransactions)
}