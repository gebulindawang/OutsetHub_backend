package middleware

import (
	"backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)
const (
	CtxUserID   = "user_id"
	CtxUsername = "username"
	CtxRole     = "role"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		const prefix = "Bearer "

		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "请先登录"})
			return
		}
		if !strings.HasPrefix(auth, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Authorization 格式应为 Bearer <token>"})
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, prefix))
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "登录已过期，请重新登录"})
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxRole, claims.Role)

		c.Next()
	}
}