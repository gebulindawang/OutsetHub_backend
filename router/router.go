package router

import (
	"backend/handler"
	"backend/middleware"

	_ "backend/docs" // 注册 swag 生成的文档,勿删

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Swagger UI:http://127.0.0.1:8080/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")

	pub := api.Group("/user")
	handler.RegisterUserRoutes(pub)

	auth := api.Group("/user")
	auth.Use(middleware.JWTAuth())
	handler.RegisterAuthRoutes(auth)

	article := api.Group("/article")
	article.Use(middleware.JWTAuth())
	handler.RegisterArticleRoutes(article)

	comment := api.Group("/comment")
	comment.Use(middleware.JWTAuth())
	handler.RegisterCommentRoutes(comment)

	category := api.Group("/category")
	category.Use(middleware.JWTAuth())
	handler.RegisterCategoryRoutes(category)

	categoryAdmin := api.Group("/category")
	categoryAdmin.Use(middleware.JWTAuth(), middleware.AdminRequired())
	handler.RegisterCategoryAdminRoutes(categoryAdmin)

	return r
}
