package router

import (
	"backend/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine{
	r := gin.Default()
	api := r.Group("/user")
	handler.RegisterUserRoutes(api)
	return  r
}