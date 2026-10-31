package handler

import (
	"backend/dao"
	"backend/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCategorys(c *gin.Context) {
	categorys, err := dao.GetCategorys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "查询成功", "data": categorys})
}

type createCategoryReq struct {
	Name string `json:"name" binding:"required"`
	Sort int    `json:"sort"`
}

func CreateCategory(c *gin.Context) {
	var req createCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误,需要 name", "err": err.Error()})
		return
	}
	if len(req.Name) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "分类名最长 50 个字符"})
		return
	}
	if err := dao.CreateCategory(model.Category{Name: req.Name, Sort: req.Sort}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "创建失败,分类名可能已存在", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "创建成功"})
}

type updateCategoryReq struct {
	Name string `json:"name"`
	Sort *int   `json:"sort"`
}

func UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的分类 ID"})
		return
	}
	var req updateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误", "err": err.Error()})
		return
	}
	updates := map[string]any{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "没有需要更新的字段"})
		return
	}
	if err := dao.UpdateCategory(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "更新失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

func DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的分类 ID"})
		return
	}
	if err := dao.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

func RegisterCategoryRoutes(rg *gin.RouterGroup) {
	rg.GET("/", GetCategorys)
}

func RegisterCategoryAdminRoutes(rg *gin.RouterGroup) {
	rg.POST("/", CreateCategory)
	rg.PUT("/:id", UpdateCategory)
	rg.DELETE("/:id", DeleteCategory)
}
