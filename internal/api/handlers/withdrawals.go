package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Xrefullx/yandexDiplom2/internal/container"
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"github.com/Xrefullx/yandexDiplom2/internal/utils"
	"github.com/Xrefullx/yandexDiplom2/internal/utils/consta"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"time"
)

func Withdraw(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	if !utils.ValidContent(c, "application/json") {
		return
	}
	storage := container.GetStorage()
	user := c.Param("loginUser")
	body, err := io.ReadAll(c.Request.Body)
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		c.String(http.StatusInternalServerError, consta.ErrorBody)
		return
	}
	var withdraw models.Withdraw
	err = json.Unmarshal(body, &withdraw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		c.String(http.StatusInternalServerError, consta.ErrorBody)
		return
	}
	withdraw.ProcessedAT, withdraw.UserLogin = time.Now(), user

	numberOrder, err := strconv.Atoi(withdraw.NumberOrder)
	if err != nil {
		c.String(http.StatusInternalServerError, "order number conversion error")
		return
	}
	if !utils.LuhValid(numberOrder) {
		c.String(http.StatusUnprocessableEntity, consta.ErrorNumberValidLuhn)
		return
	}
	err = storage.AddWithdraw(ctx, withdraw)
	if err != nil {
		if errors.Is(err, consta.ErrorStatusShortfallAccount) {
			c.String(http.StatusPaymentRequired, consta.ErrorStatusShortfallAccount.Error())
			return
		}
		c.String(http.StatusInternalServerError, consta.ErrorDataBase)
		return
	}
	c.String(http.StatusOK, "the write-off has been completed")
}
func GetWithdrawsByUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	storage := container.GetStorage()
	user := c.Param("loginUser")
	orders, err := storage.GetWithdraws(ctx, user)
	if err != nil {
		c.String(http.StatusInternalServerError, consta.ErrorDataBase)
		return
	}
	if len(orders) == 0 {
		c.String(http.StatusNoContent, "no data")
		return
	}
	c.JSON(http.StatusOK, orders)
}
