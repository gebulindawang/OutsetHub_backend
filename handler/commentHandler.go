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

// GetComments 评论列表
//
//	@Summary		某篇文章的评论列表
//	@Description	按 article_id 查询,只返回 status=1 的评论,按创建时间升序
//	@Tags			comment
//	@Produce		json
//	@Security		BearerAuth
//	@Param			article_id	query		int	true	"文章 ID"
//	@Success		200			{object}	model.CommentListResp
//	@Failure		400			{object}	model.ErrResp		"缺少 article_id"
//	@Failure		401			{object}	model.CodeMsgResp	"未登录或 token 失效"
//	@Failure		500			{object}	model.ErrResp
//	@Router			/comment/ [get]
func GetComments(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Query("article_id"), 10, 64)
	if err != nil || articleID == 0 {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "需要 article_id 参数"})
		return
	}
	comments, err := dao.GetComments(articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "查询失败", Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.CommentListResp{Msg: "查询成功", Data: comments})
}

// GetCommentByID 评论详情
//
//	@Summary		评论详情
//	@Description	按 ID 查询单条评论
//	@Tags			comment
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"评论 ID"
//	@Success		200	{object}	model.CommentResp
//	@Failure		400	{object}	model.ErrResp		"无效的评论 ID"
//	@Failure		401	{object}	model.CodeMsgResp	"未登录或 token 失效"
//	@Failure		404	{object}	model.ErrResp		"评论不存在"
//	@Failure		500	{object}	model.ErrResp
//	@Router			/comment/{id} [get]
func GetCommentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "无效的评论 ID"})
		return
	}
	comment, err := dao.GetCommentByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, model.ErrResp{Msg: "评论不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "查询失败", Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.CommentResp{Msg: "查询成功", Data: comment})
}

// CreateComment 发表评论
//
//	@Summary		发表评论
//	@Description	评论人取自 JWT。parent_id 为 0 表示一级评论,非 0 表示回复某条评论。正文最长 500 字。
//	@Tags			comment
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		model.CreateCommentReq	true	"评论内容"
//	@Success		200		{object}	model.MsgResp
//	@Failure		400		{object}	model.ErrResp		"缺少 article_id/content 或超出长度"
//	@Failure		401		{object}	model.CodeMsgResp	"未登录或 token 失效"
//	@Failure		500		{object}	model.ErrResp
//	@Router			/comment/ [post]
func CreateComment(c *gin.Context) {
	var req model.CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "参数错误,需要 article_id/content", Err: err.Error()})
		return
	}
	if len(req.Content) > 500 {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "评论最长 500 个字符"})
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
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "创建失败", Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.MsgResp{Msg: "发布成功"})
}

// UpdateComment 修改评论
//
//	@Summary		修改评论
//	@Description	仅评论作者本人可操作,正文不能为空且最长 500 字
//	@Tags			comment
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"评论 ID"
//	@Param			body	body		model.UpdateCommentReq	true	"新内容"
//	@Success		200		{object}	model.MsgResp
//	@Failure		400		{object}	model.ErrResp		"内容为空或超出长度"
//	@Failure		401		{object}	model.CodeMsgResp	"未登录或 token 失效"
//	@Failure		403		{object}	model.ErrResp		"只能编辑自己的评论"
//	@Failure		404		{object}	model.ErrResp		"评论不存在"
//	@Failure		500		{object}	model.ErrResp
//	@Router			/comment/{id} [put]
func UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "无效的评论 ID"})
		return
	}
	comment, err := dao.GetCommentByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, model.ErrResp{Msg: "评论不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "查询失败", Err: err.Error()})
		return
	}
	if comment.UserID != c.GetInt(middleware.CtxUserID) {
		c.JSON(http.StatusForbidden, model.ErrResp{Msg: "只能编辑自己的评论"})
		return
	}

	var req model.UpdateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "参数错误", Err: err.Error()})
		return
	}
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "没有需要更新的字段"})
		return
	}
	if len(req.Content) > 500 {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "评论最长 500 个字符"})
		return
	}
	updates := map[string]any{"content": req.Content}
	if err := dao.UpdateComment(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "更新失败", Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.MsgResp{Msg: "修改成功"})
}

// DeleteComment 删除评论
//
//	@Summary		删除评论
//	@Description	软删除,status 置为 2。仅评论作者本人可操作。
//	@Tags			comment
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"评论 ID"
//	@Success		200	{object}	model.MsgResp
//	@Failure		400	{object}	model.ErrResp		"无效的评论 ID"
//	@Failure		401	{object}	model.CodeMsgResp	"未登录或 token 失效"
//	@Failure		403	{object}	model.ErrResp		"只能删除自己的评论"
//	@Failure		404	{object}	model.ErrResp		"评论不存在"
//	@Failure		500	{object}	model.ErrResp
//	@Router			/comment/{id} [delete]
func DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "无效的评论 ID"})
		return
	}
	comment, err := dao.GetCommentByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, model.ErrResp{Msg: "评论不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "查询失败", Err: err.Error()})
		return
	}
	if comment.UserID != c.GetInt(middleware.CtxUserID) {
		c.JSON(http.StatusForbidden, model.ErrResp{Msg: "只能删除自己的评论"})
		return
	}
	if err := dao.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "删除失败", Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.MsgResp{Msg: "删除成功"})
}

func RegisterCommentRoutes(rg *gin.RouterGroup) {
	rg.GET("/", GetComments)
	rg.GET("/:id", GetCommentByID)
	rg.POST("/", CreateComment)
	rg.PUT("/:id", UpdateComment)
	rg.DELETE("/:id", DeleteComment)
}
