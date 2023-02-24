package handlers

import (
	"context"
	"github.com/Xrefullx/yandexDiplom2/internal/container"
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"github.com/Xrefullx/yandexDiplom2/internal/utils"
	"github.com/Xrefullx/yandexDiplom2/internal/utils/consta"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	if !utils.ValidContent(c, "application/json") {
		return
	}
	storage := container.GetStorage()
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		c.String(http.StatusInternalServerError, consta.ErrorBody)
		return
	}
	if user.Login == "" || user.Password == "" {
		c.String(http.StatusBadRequest, "login or pass is not valid")
		return
	}
	authenticationUser, err := storage.Authentication(ctx, user)
	if err != nil {
		c.String(http.StatusInternalServerError, consta.ErrorDataBase)
		return
	}
	if !authenticationUser {
		c.String(http.StatusUnauthorized, "login or pass is not correct")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 100)),
			IssuedAt:  jwt.At(time.Now())},
		Login: user.Login,
	})
	accessToken, err := token.SignedString([]byte(container.GetConfig().SecretKey))
	if err != nil {
		c.String(http.StatusInternalServerError, "error generation jwt")
		return
	}
	c.Header("Authorization", "Bearer "+accessToken)
}
