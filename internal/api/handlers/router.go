package handlers

import (
	"github.com/Xrefullx/yandexDiplom2/internal/api/middleware"
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"github.com/gin-gonic/gin"
)

func Router(cfg models.Config) *gin.Engine {
	if cfg.ReleaseMOD {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.JwtValid())

	gUser := r.Group("/api/user")
	{
		gUser.POST("/register", Register)
		gUser.POST("/login", Login)
		gUser.POST("/orders", AddOrders)
		gUser.GET("/orders", GetOrders)
		gUser.GET("/balance", UserBalance)
		gUser.POST("/balance/withdraw", AddWithdraw)
		gUser.GET("/withdrawals", GetWithdraws)
	}
	return r
}
