package middleware

import (
	"backend/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetInt(CtxRole) != 1 {
			response.Abort(c, http.StatusForbidden, "需要管理员权限")
			return
		}
		c.Next()
	}
}
