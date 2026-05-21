package delivery

import "github.com/gin-gonic/gin"

func RegisterFnBRoutes(rg *gin.RouterGroup, menuH *MenuHandler, orderH *OrderHandler) {
	// Master Product
	rg.POST("/products", menuH.CreateProduct)
	rg.GET("/products", menuH.GetProducts)
	rg.PUT("/products/:id", menuH.UpdateProduct)
	rg.DELETE("/products/:id", menuH.DeleteProduct)
	rg.PUT("/products/:id/toggle", menuH.ToggleAvailability)

	// Order & Dapur (SEKARANG SUDAH AKTIF)
	rg.POST("/order", orderH.CreateOrder)
	rg.GET("/kitchen", orderH.GetAntreanDapur)
	rg.PUT("/kitchen/:id/selesai", orderH.SelesaikanOrderan)
}