package router

import (
	"github.com/Xrefullx/yandexDiplom2/internal/api/handlers"
	"github.com/Xrefullx/yandexDiplom2/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.JwtValid())

	gUser := r.Group("/api/user")
	{
		gUser.POST("/register", handlers.RegisterUser)
		gUser.POST("/login", handlers.Login)
		gUser.POST("/orders", handlers.AddOrdersByUser)
		gUser.GET("/orders", handlers.GetOrdersByUser)
		gUser.GET("/balance", handlers.UserBalance)
		gUser.POST("/balance/withdraw", handlers.Withdraw)
		gUser.GET("/withdrawals", handlers.GetWithdrawsByUser)
	}
	return r
}
