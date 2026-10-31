package handler

import (
	"backend/dao"
	"backend/middleware"
	"backend/model"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetArticles(c *gin.Context) {
	articles, err := dao.GetArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "查询成功", "data": articles})
}

func GetArticleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的文章 ID"})
		return
	}
	article, err := dao.GetArticleByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"msg": "文章不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "查询成功", "data": article})
}

type createArticleReq struct {
	Title      string `json:"title" binding:"required"`
	Summary    string `json:"summary"`
	Content    string `json:"content" binding:"required"`
	CategoryID int    `json:"category_id" binding:"required"`
	Status     *int   `json:"status"`
}

func CreateArticle(c *gin.Context) {
	var req createArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误,需要 title/content/category_id", "err": err.Error()})
		return
	}

	article := model.Article{
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		UserID:     c.GetInt(middleware.CtxUserID),
		CategoryID: req.CategoryID,
		Status:     1,
	}
	if req.Status != nil {
		article.Status = *req.Status
	}
	if err := dao.CreateArticle(article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "发布成功"})
}

type updateArticleReq struct {
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Content    string `json:"content"`
	CategoryID *int   `json:"category_id"`
	Status     *int   `json:"status"`
}

func UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的文章 ID"})
		return
	}
	article, err := dao.GetArticleByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"msg": "文章不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败", "err": err.Error()})
		return
	}
	if article.UserID != c.GetInt(middleware.CtxUserID) {
		c.JSON(http.StatusForbidden, gin.H{"msg": "只能编辑自己的文章"})
		return
	}

	var req updateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误", "err": err.Error()})
		return
	}
	updates := map[string]any{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "没有需要更新的字段"})
		return
	}
	if err := dao.UpdateArticle(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "更新失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "修改成功"})
}

func DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的文章 ID"})
		return
	}
	article, err := dao.GetArticleByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"msg": "文章不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败", "err": err.Error()})
		return
	}
	if article.UserID != c.GetInt(middleware.CtxUserID) {
		c.JSON(http.StatusForbidden, gin.H{"msg": "只能删除自己的文章"})
		return
	}
	if err := dao.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

func RegisterArticleRoutes(rg *gin.RouterGroup) {
	rg.GET("/", GetArticles)
	rg.GET("/:id", GetArticleByID)
	rg.POST("/", CreateArticle)
	rg.PUT("/:id", UpdateArticle)
	rg.DELETE("/:id", DeleteArticle)
}
