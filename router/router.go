package router

import (
	"backend/handler"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine{
	r := gin.Default()
	api := r.Group("/api")


	pub := api.Group("/user")
	handler.RegisterUserRoutes(pub)

	auth := api.Group("/user")
	auth.Use(middleware.JWTAuth())
	handler.RegisterAuthRoutes(auth)

	article := api.Group("/article")
	article.Use(middleware.JWTAuth())
	handler.RegisterArticleRoutes(article)

	category := api.Group("/category")
	category.Use(middleware.JWTAuth())
	handler.RegisterCategoryRoutes(category)
	return  r
}