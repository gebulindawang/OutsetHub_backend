package handler

import (
	"backend/dao"
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCategorys(c *gin.Context) {
	var categorys []model.Category
	categorys,err := dao.GetCategorys()
	if err != nil{
		panic(err)
	}
	c.JSON(http.StatusOK,gin.H{"msg":"操作成功","data":categorys})
}

func RegisterCategoryRoutes(rg *gin.RouterGroup){
	rg.GET("/",GetArticles)
}