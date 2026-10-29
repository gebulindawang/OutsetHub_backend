package router

import (
	"backend/handler"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine{
	r := gin.Default()
	api := r.Group("/user")
	handler.RegisterUserRoutes(api)

	auth := r.Group("/user")
	auth.Use(middleware.JWTAuth())
	handler.RegisterAuthRoutes(auth)

	article := r.Group("/article")
	article.Use(middleware.JWTAuth())
	handler.RegisterArticleRoutes(article)
	return  r
}