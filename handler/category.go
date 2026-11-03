package handler

import (
	"backend/dao"
	"backend/model"
	"backend/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCategorys 板块列表
//
//	@Summary		板块列表
//	@Description	按 sort 升序返回全部板块,登录用户即可读取
//	@Tags			category
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	model.CategoryListResp
//	@Failure		401	{object}	model.ErrResp	"未登录或 token 失效"
//	@Failure		500	{object}	model.ErrResp
//	@Router			/category/ [get]
func GetCategorys(c *gin.Context) {
	categorys, err := dao.GetCategorys()
	if err != nil {
		response.FailErr(c, http.StatusInternalServerError, "查询失败", err)
		return
	}
	response.Ok(c, categorys)
}

// CreateCategory 新建板块
//
//	@Summary		新建板块
//	@Description	仅管理员(role=1)可调用。板块名全局唯一,最长 50 个字符。
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		model.CreateCategoryReq	true	"板块信息"
//	@Success		200		{object}	model.MsgResp
//	@Failure		400		{object}	model.ErrResp	"参数错误或板块名已存在"
//	@Failure		401		{object}	model.ErrResp	"未登录或 token 失效"
//	@Failure		403		{object}	model.ErrResp	"需要管理员权限"
//	@Router			/category/ [post]
func CreateCategory(c *gin.Context) {
	var req model.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailErr(c, http.StatusBadRequest, "参数错误,需要 name", err)
		return
	}
	if len(req.Name) > 50 {
		response.Fail(c, http.StatusBadRequest, "分类名最长 50 个字符")
		return
	}
	if err := dao.CreateCategory(model.Category{Name: req.Name, Sort: req.Sort}); err != nil {
		response.FailErr(c, http.StatusBadRequest, "创建失败,分类名可能已存在", err)
		return
	}
	response.OkMsg(c, "创建成功")
}

// UpdateCategory 修改板块
//
//	@Summary		修改板块
//	@Description	仅管理员可调用,部分更新:只传需要修改的字段
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"板块 ID"
//	@Param			body	body		model.UpdateCategoryReq	true	"待更新字段"
//	@Success		200		{object}	model.MsgResp
//	@Failure		400		{object}	model.ErrResp	"无效的 ID 或没有可更新的字段"
//	@Failure		401		{object}	model.ErrResp	"未登录或 token 失效"
//	@Failure		403		{object}	model.ErrResp	"需要管理员权限"
//	@Failure		500		{object}	model.ErrResp
//	@Router			/category/{id} [put]
func UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, "无效的分类 ID")
		return
	}
	var req model.UpdateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailErr(c, http.StatusBadRequest, "参数错误", err)
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
		response.Fail(c, http.StatusBadRequest, "没有需要更新的字段")
		return
	}
	if err := dao.UpdateCategory(id, updates); err != nil {
		response.FailErr(c, http.StatusInternalServerError, "更新失败", err)
		return
	}
	response.OkMsg(c, "更新成功")
}

// DeleteCategory 删除板块
//
//	@Summary		删除板块
//	@Description	仅管理员可调用。物理删除,不校验是否仍有文章引用该板块。
//	@Tags			category
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"板块 ID"
//	@Success		200	{object}	model.MsgResp
//	@Failure		400	{object}	model.ErrResp	"无效的分类 ID"
//	@Failure		401	{object}	model.ErrResp	"未登录或 token 失效"
//	@Failure		403	{object}	model.ErrResp	"需要管理员权限"
//	@Failure		500	{object}	model.ErrResp
//	@Router			/category/{id} [delete]
func DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, "无效的分类 ID")
		return
	}
	if err := dao.DeleteCategory(id); err != nil {
		response.FailErr(c, http.StatusInternalServerError, "删除失败", err)
		return
	}
	response.OkMsg(c, "删除成功")
}

func RegisterCategoryRoutes(rg *gin.RouterGroup) {
	rg.GET("/", GetCategorys)
}

func RegisterCategoryAdminRoutes(rg *gin.RouterGroup) {
	rg.POST("/", CreateCategory)
	rg.PUT("/:id", UpdateCategory)
	rg.DELETE("/:id", DeleteCategory)
}
