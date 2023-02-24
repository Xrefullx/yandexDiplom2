package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Xrefullx/yandexDiplom2/internal/container"
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"github.com/Xrefullx/yandexDiplom2/internal/utils"
	"github.com/Xrefullx/yandexDiplom2/internal/utils/consta"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func AddOrdersByUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	if !utils.ValidContent(c, "text/plain") {
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
	var numberOrder int
	err = json.Unmarshal(body, &numberOrder)
	if err != nil {
		c.String(http.StatusInternalServerError, consta.ErrorBody)
		return
	}
	if !utils.LuhValid(numberOrder) {
		c.String(http.StatusUnprocessableEntity, consta.ErrorNumberValidLuhn)
		return
	}
	numberOrderStr := fmt.Sprint(numberOrder)
	err = storage.AddOrder(ctx, numberOrderStr,
		models.Order{
			NumberOrder: numberOrderStr,
			UserLogin:   user,
			Status:      consta.OrderStatusNEW,
			Uploaded:    time.Now(),
		})
	if err != nil {
		if errors.Is(err, consta.ErrorNoUNIQUE) {
			order, errGet := storage.GetOrder(ctx, numberOrderStr)
			if errGet != nil {
				c.String(http.StatusInternalServerError, consta.ErrorDataBase)
				return
			}
			if order.UserLogin == user {
				c.String(http.StatusOK, "the order number has already been uploaded by this user")
				return
			}
			c.String(http.StatusConflict, "the order number has already been uploaded by this user")
			return
		}
		c.String(http.StatusInternalServerError, consta.ErrorDataBase)
		return
	}
	c.String(http.StatusAccepted, "the new order number has been accepted for processing")
}
func GetOrdersByUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	storage := container.GetStorage()
	user := c.Param("loginUser")
	orders, err := storage.GetOrders(ctx, user)
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
