package api

import (
	_ "main/api/docs"
	"main/api/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func NewServer(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.POST("/branches", h.CreateBranch)
	r.GET("/branches", h.GetAllBranch)
	r.GET("/branches/:id", h.GetBranch)
	r.PUT("/branches/:id", h.UpdateBranch)
	r.DELETE("/branches/:id", h.DeleteBranch)

	// r.POST("/sales", h.CreateSale)
	// r.GET("/sales", h.GetAllSale)
	// r.GET("/sales/:id", h.GetSale)
	// r.PUT("/sales/:id", h.UpdateSale)
	// r.DELETE("/sales/:id", h.DeleteSale)

	r.POST("/staffs", h.CreateStaff)
	r.GET("/staffs", h.GetAllStaff)
	r.GET("/staffs/:id", h.GetStaff)
	r.PUT("/staffs/:id", h.UpdateStaff)
	r.DELETE("/staffs/:id", h.DeleteStaff)

	r.POST("/tarifs", h.CreateStaffTarif)
	r.GET("/tarifs", h.GetAllStaffTarif)
	r.GET("/tarifs/:id", h.GetStaffTarif)
	r.PUT("/tarifs/:id", h.UpdateStaffTarif)
	r.DELETE("/tarifs/:id", h.DeleteStaffTarif)

	// r.POST("/transactions", h.CreateStaffTransaction)
	// r.GET("/transactions", h.GetAllStaffTransaction)
	// r.GET("/transactions/:id", h.GetStaffTransaction)
	// r.PUT("/transactions/:id", h.UpdateStaffTransaction)
	// r.DELETE("/transactions/:id", h.DeleteStaffTransaction)

	r.PUT("/change-password/:id", h.ChangePassword)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return r
}
