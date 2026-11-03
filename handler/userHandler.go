package handler

import (
	"backend/dao"
	"backend/middleware"
	"backend/model"
	"backend/response"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register 注册新用户
//
//	@Summary		用户注册
//	@Description	用户名全局唯一,密码经 bcrypt 哈希后入库。成功仅返回提示信息,不返回 token。
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.RegisterReq	true	"注册参数"
//	@Success		200		{object}	model.MsgResp
//	@Failure		400		{object}	model.ErrResp	"参数错误或用户名已存在"
//	@Router			/user/register [post]
func Register(c *gin.Context) {
	var req model.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailErr(c, http.StatusBadRequest, "参数校验失败", err)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.FailErr(c, http.StatusBadRequest, "注册失败", err)
		return
	}
	user := model.User{
		Username: req.Username,
		Password: string(hash),
		Nickname: req.Nickname,
	}
	if err := dao.Register(&user); err != nil {
		response.FailErr(c, http.StatusBadRequest, "注册失败,用户名可能已存在", err)
		return
	}
	response.OkMsg(c, "注册成功")
}

// Login 登录并签发 JWT
//
//	@Summary		用户登录
//	@Description	校验用户名密码,返回有效期 2 小时的 JWT。token 位于 data.token,后续请求放在 Authorization: Bearer <token>
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.LoginReq	true	"登录参数"
//	@Success		200		{object}	model.LoginResp
//	@Failure		400		{object}	model.ErrResp	"用户名或密码错误"
//	@Router			/user/login [post]
func Login(c *gin.Context) {
	var req model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailErr(c, http.StatusBadRequest, "数据解析失败", err)
		return
	}
	u, err := dao.FindUserByusername(req.Username)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "用户名或密码错误")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		response.Fail(c, http.StatusBadRequest, "用户名或密码错误")
		return
	}
	token, err := utils.GenToken(u.ID, u.Username, u.Role)
	if err != nil {
		response.FailErr(c, http.StatusInternalServerError, "签发 token 失败", err)
		return
	}
	response.Ok(c, model.LoginData{Token: token})
}

// UserInfo 获取当前登录用户
//
//	@Summary		当前登录用户信息
//	@Description	从 JWT 中取用户名再查库,Password 字段不会返回
//	@Tags			user
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	model.UserResp
//	@Failure		401	{object}	model.ErrResp	"未登录或 token 失效"
//	@Failure		500	{object}	model.ErrResp
//	@Router			/user/me [get]
func UserInfo(c *gin.Context) {
	username := c.GetString(middleware.CtxUsername)
	u, err := dao.FindUserByusername(username)
	if err != nil {
		response.FailErr(c, http.StatusInternalServerError, "查询失败", err)
		return
	}
	response.Ok(c, u)
}

func RegisterUserRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", Register)
	rg.POST("/login", Login)
}

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	rg.GET("/me", UserInfo)
}
