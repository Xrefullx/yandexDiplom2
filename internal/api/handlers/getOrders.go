package handlers

import (
	"context"
	"github.com/Xrefullx/yandexDiplom2/internal/api/consta"
	"github.com/Xrefullx/yandexDiplom2/internal/api/container"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func GetOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	log := container.GetLog()
	storage := container.GetStorage()
	user := c.Param("loginUser")
	log.Debug("received a request to display orders",
		zap.String("loginUser", user))

	orders, err := storage.GetOrders(ctx, user)
	if err != nil {
		log.Error(consta.ErrorDataBase, zap.Error(err), zap.String("func", "GetOrders"))
		c.String(http.StatusInternalServerError, consta.ErrorDataBase)
		return
	}
	if len(orders) == 0 {
		log.Debug("no data for request", zap.String("loginUser", user))
		c.String(http.StatusNoContent, "no data for request")
		return
	}
	c.JSON(http.StatusOK, orders)
}
