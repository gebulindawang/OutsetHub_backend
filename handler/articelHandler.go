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

// GetArticles 文章列表
//
//	@Summary		文章列表
//	@Description	只返回 status=1 的已发布文章,按创建时间倒序,当前固定 10 条,暂不支持分页与按板块筛选。
//	@Tags			article
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	model.ArticleListResp
//	@Failure		401	{object}	model.CodeMsgResp	"未登录或 token 失效"
//	@Failure		500	{object}	model.ErrResp
//	@Router			/article/ [get]
func GetArticles(c *gin.Context) {
	articles, err := dao.GetArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "查询失败", Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ArticleListResp{Msg: "查询成功", Data: articles})
}

// GetArticleByID 文章详情
//
//	@Summary		文章详情
//	@Description	按 ID 查询单篇文章。注意:本接口不校验 status,草稿和已删除的文章也能取到,由调用方自行判断。
//	@Tags			article
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"文章 ID"
//	@Success		200	{object}	model.ArticleResp
//	@Failure		400	{object}	model.ErrResp		"无效的文章 ID"
//	@Failure		401	{object}	model.CodeMsgResp	"未登录或 token 失效"
//	@Failure		404	{object}	model.ErrResp		"文章不存在"
//	@Failure		500	{object}	model.ErrResp
//	@Router			/article/{id} [get]
func GetArticleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "无效的文章 ID"})
		return
	}
	article, err := dao.GetArticleByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, model.ErrResp{Msg: "文章不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "查询失败", Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ArticleResp{Msg: "查询成功", Data: article})
}

// CreateArticle 发布文章
//
//	@Summary		发布文章
//	@Description	作者取自 JWT,不信任请求体里的 user_id。status 不传默认 1 已发布,传 0 表示存草稿。
//	@Tags			article
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		model.CreateArticleReq	true	"文章内容"
//	@Success		200		{object}	model.MsgResp
//	@Failure		400		{object}	model.ErrResp		"缺少 title/content/category_id"
//	@Failure		401		{object}	model.CodeMsgResp	"未登录或 token 失效"
//	@Failure		500		{object}	model.ErrResp
//	@Router			/article/ [post]
func CreateArticle(c *gin.Context) {
	var req model.CreateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "参数错误,需要 title/content/category_id", Err: err.Error()})
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
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "创建失败", Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.MsgResp{Msg: "发布成功"})
}

// UpdateArticle 编辑文章
//
//	@Summary		编辑文章
//	@Description	仅作者本人可操作。部分更新:只提交需要修改的字段,空字符串视为不修改。
//	@Tags			article
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"文章 ID"
//	@Param			body	body		model.UpdateArticleReq	true	"待更新字段"
//	@Success		200		{object}	model.MsgResp
//	@Failure		400		{object}	model.ErrResp		"参数错误或没有可更新的字段"
//	@Failure		401		{object}	model.CodeMsgResp	"未登录或 token 失效"
//	@Failure		403		{object}	model.ErrResp		"只能编辑自己的文章"
//	@Failure		404		{object}	model.ErrResp		"文章不存在"
//	@Failure		500		{object}	model.ErrResp
//	@Router			/article/{id} [put]
func UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "无效的文章 ID"})
		return
	}
	article, err := dao.GetArticleByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, model.ErrResp{Msg: "文章不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "查询失败", Err: err.Error()})
		return
	}
	if article.UserID != c.GetInt(middleware.CtxUserID) {
		c.JSON(http.StatusForbidden, model.ErrResp{Msg: "只能编辑自己的文章"})
		return
	}

	var req model.UpdateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "参数错误", Err: err.Error()})
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
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "没有需要更新的字段"})
		return
	}
	if err := dao.UpdateArticle(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "更新失败", Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.MsgResp{Msg: "修改成功"})
}

// DeleteArticle 删除文章
//
//	@Summary		删除文章
//	@Description	软删除,把 status 置为 2,记录仍留在库中。仅作者本人可操作。
//	@Tags			article
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"文章 ID"
//	@Success		200	{object}	model.MsgResp
//	@Failure		400	{object}	model.ErrResp		"无效的文章 ID"
//	@Failure		401	{object}	model.CodeMsgResp	"未登录或 token 失效"
//	@Failure		403	{object}	model.ErrResp		"只能删除自己的文章"
//	@Failure		404	{object}	model.ErrResp		"文章不存在"
//	@Failure		500	{object}	model.ErrResp
//	@Router			/article/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, model.ErrResp{Msg: "无效的文章 ID"})
		return
	}
	article, err := dao.GetArticleByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, model.ErrResp{Msg: "文章不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "查询失败", Err: err.Error()})
		return
	}
	if article.UserID != c.GetInt(middleware.CtxUserID) {
		c.JSON(http.StatusForbidden, model.ErrResp{Msg: "只能删除自己的文章"})
		return
	}
	if err := dao.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrResp{Msg: "删除失败", Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.MsgResp{Msg: "删除成功"})
}

func RegisterArticleRoutes(rg *gin.RouterGroup) {
	rg.GET("/", GetArticles)
	rg.GET("/:id", GetArticleByID)
	rg.POST("/", CreateArticle)
	rg.PUT("/:id", UpdateArticle)
	rg.DELETE("/:id", DeleteArticle)
}
