package middleware

import (
	"github.com/Xrefullx/yandexDiplom2/internal/container"
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/auth/pkg/auth"
	"net/http"
	"strings"
)

func JwtValid() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/user/register" || c.Request.URL.Path == "/api/user/login" {
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		header := strings.Split(authHeader, " ")
		if len(header) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if header[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		login, err := parseToken(header[1],
			[]byte(container.Container.Get("server-config").(models.Config).SecretKey),
		)
		if err != nil {
			status := http.StatusBadRequest
			if err == auth.ErrInvalidAccessToken {
				status = http.StatusUnauthorized
			}

			c.AbortWithStatus(status)
			return
		}
		c.AddParam("loginUser", login)
	}
}
