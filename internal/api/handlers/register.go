package handlers

import (
	"context"
	"github.com/Xrefullx/yandexDiplom2/internal/container"
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"github.com/Xrefullx/yandexDiplom2/internal/utils"
	"github.com/Xrefullx/yandexDiplom2/internal/utils/consta"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	if !utils.ValidContent(c, "application/json") {
		return
	}
	storage := container.GetStorage()
	var user models.User
	// Проверяем правильность формата запроса
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		c.String(http.StatusInternalServerError, consta.ErrorBody)
		return
	}
	if user.Login == "" || user.Password == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err := storage.Adduser(ctx, user)
	if err != nil {
		if errors.Is(err, consta.ErrorNoUNIQUE) {
			{
				c.JSON(http.StatusConflict, gin.H{"error": "Login already exists"})
				return
			}
			c.String(http.StatusInternalServerError, consta.ErrorDataBase)
			return
		}
	}
	authToken := "xrefullxAuth"
	// Устанавливаем куку с токеном аутентификации
	c.SetCookie("auth_token", authToken, 60*60*24, "/", "localhost", false, true)
	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "User registered and authenticated"})
	c.Redirect(http.StatusPermanentRedirect, "/api/user/login")
}
