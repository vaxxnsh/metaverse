package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaxxnsh/metaverse/api/internal/lib/response"
	"github.com/vaxxnsh/metaverse/api/internal/utils"
)

const (
	UserIDKey  = "userId"
	AdminIDKey = "adminId"
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			response.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing token", struct{}{})
			c.Abort()
			return
		}

		userId, err := utils.ParseToken(token)
		if err != nil {
			response.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token", struct{}{})
			c.Abort()
			return
		}

		c.Set(UserIDKey, userId)
		c.Next()
	}
}

func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			response.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing token", struct{}{})
			c.Abort()
			return
		}

		adminId, err := utils.ParseToken(token)
		if err != nil {
			response.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token", struct{}{})
			c.Abort()
			return
		}

		c.Set(AdminIDKey, adminId)
		c.Next()
	}
}
