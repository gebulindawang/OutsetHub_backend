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

func GetComments(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Query("article_id"), 10, 64)
	if err != nil || articleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "需要 article_id 参数"})
		return
	}
	comments, err := dao.GetComments(articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "查询成功", "data": comments})
}

func GetCommentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的评论 ID"})
		return
	}
	comment, err := dao.GetCommentByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"msg": "评论不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "查询成功", "data": comment})
}

type createCommentReq struct {
	ArticleID     uint64 `json:"article_id" binding:"required"`
	ParentID      uint64 `json:"parent_id"`
	ReplyToUserID int    `json:"reply_to_user_id"`
	Content       string `json:"content" binding:"required"`
}

func CreateComment(c *gin.Context) {
	var req createCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误,需要 article_id/content", "err": err.Error()})
		return
	}
	if len(req.Content) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "评论最长 500 个字符"})
		return
	}
	comment := model.Comment{
		ArticleID:     req.ArticleID,
		UserID:        c.GetInt(middleware.CtxUserID),
		ParentID:      req.ParentID,
		ReplyToUserID: req.ReplyToUserID,
		Content:       req.Content,
		Status:        1,
	}
	if err := dao.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "发布成功"})
}

type updateCommentReq struct {
	Content string `json:"content"`
}

func UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的评论 ID"})
		return
	}
	comment, err := dao.GetCommentByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"msg": "评论不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败", "err": err.Error()})
		return
	}
	if comment.UserID != c.GetInt(middleware.CtxUserID) {
		c.JSON(http.StatusForbidden, gin.H{"msg": "只能编辑自己的评论"})
		return
	}

	var req updateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误", "err": err.Error()})
		return
	}
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "没有需要更新的字段"})
		return
	}
	if len(req.Content) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "评论最长 500 个字符"})
		return
	}
	updates := map[string]any{"content": req.Content}
	if err := dao.UpdateComment(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "更新失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "修改成功"})
}

func DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的评论 ID"})
		return
	}
	comment, err := dao.GetCommentByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"msg": "评论不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败", "err": err.Error()})
		return
	}
	if comment.UserID != c.GetInt(middleware.CtxUserID) {
		c.JSON(http.StatusForbidden, gin.H{"msg": "只能删除自己的评论"})
		return
	}
	if err := dao.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

func RegisterCommentRoutes(rg *gin.RouterGroup) {
	rg.GET("/", GetComments)
	rg.GET("/:id", GetCommentByID)
	rg.POST("/", CreateComment)
	rg.PUT("/:id", UpdateComment)
	rg.DELETE("/:id", DeleteComment)
}
