package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetInt(CtxRole) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "msg": "需要管理员权限"})
			return
		}
		c.Next()
	}
}
