package handler

import (
	"backend/dao"
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetArticles(c *gin.Context) {
	 res,err := dao.GetArticles()
	 if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"err":err.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{"msg":"查询成功","data":res})
}

func CreateArticle(c *gin.Context){
	var article model.Article
	if err := c.ShouldBindJSON(&article);err != nil{
		panic(err)
	}
	err := dao.CreateArticle(article)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"err":err.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{"msg":"操作成功"})
}

func DeleteArticle(c *gin.Context){
	var article model.Article
	if err := c.ShouldBindJSON(&article);err!=nil{
		panic(err)
	}
	err := dao.DeleteArticle(article)
	if err != nil{
		panic(err)
	}
	c.JSON(http.StatusOK,gin.H{"msg":"删除成功"})
}

func UpdateArticle(c *gin.Context){
	var article model.Article
	if err := c.ShouldBindJSON(&article);err != nil{
		panic(err)
	}
	err := dao.UpdateArticle(article)
	if err != nil{
		panic(err)
	}
	c.JSON(http.StatusOK,gin.H{"msg":"修改成功"})
}
func RegisterArticleRoutes(rg *gin.RouterGroup){
	rg.GET("/",GetArticles)
	rg.POST("/",CreateArticle)
	rg.DELETE("/",DeleteArticle)
	rg.PUT("/",UpdateArticle)
}