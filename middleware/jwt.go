package middleware

import (
	"backend/response"
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
			response.Abort(c, http.StatusUnauthorized, "请先登录")
			return
		}
		if !strings.HasPrefix(auth, prefix) {
			response.Abort(c, http.StatusUnauthorized, "Authorization 格式应为 Bearer <token>")
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, prefix))
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			response.Abort(c, http.StatusUnauthorized, "登录已过期,请重新登录")
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxRole, claims.Role)

		c.Next()
	}
}
