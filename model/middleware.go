package model

import (
	"net/http"
	"strings"

	"github.com/EnTing0417/go-lib/model"
	"github.com/gin-gonic/gin"
)

func containsRoute(routes []string, target string) bool {
	for _, route := range routes {
		if strings.Contains(target, route) {
			return true
		}
	}
	return false
}

func AuthMiddleware(protectedRoutes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if containsRoute(protectedRoutes, c.Request.URL.Path) {
			token := c.GetHeader("Authorization")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
				c.Abort()
				return
			}
			config := model.ReadConfig()
			token = token[7:len(token)]

			publicKey, err := model.ParseRSAPublicKeyFromConfig(config.Auth.TkPublicKey)

			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid public key"})
				c.Abort()
				return
			}

			claims, isValidToken := model.IsTokenValid(token, publicKey)

			if !isValidToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}
			c.Set("claims", claims)
		}
		c.Next()
	}
}
