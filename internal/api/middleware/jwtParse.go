package middleware

import (
	"fmt"
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/zhashkevych/auth/pkg/auth"
)

func parseToken(accessToken string, signingKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims.Login, nil
	}

	return "", auth.ErrInvalidAccessToken
}
