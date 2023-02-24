package handlers

import (
	"context"
	"github.com/Xrefullx/yandexDiplom2/internal/container"
	"github.com/Xrefullx/yandexDiplom2/internal/utils/consta"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserBalance(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	storage := container.GetStorage()
	user := c.Param("loginUser")
	sum, spent, err := storage.GetUserBalance(ctx, user)
	if err != nil {
		c.String(http.StatusInternalServerError, consta.ErrorDataBase)
		return
	}
	response := map[string]float64{
		"current":   sum - spent,
		"withdrawn": spent,
	}
	c.JSONP(http.StatusOK, response)
}
