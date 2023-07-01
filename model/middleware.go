package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"github.com/EnTing0417/go-lib/model"
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
        
			claims, isValidToken := model.IsTokenValid(token, config.Auth.SecretKey)
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
