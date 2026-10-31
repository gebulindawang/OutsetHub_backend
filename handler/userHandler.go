package handler

import (
	"backend/dao"
	"backend/middleware"
	"backend/model"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registerReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
}

func Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	user := model.User{
		Username: req.Username,
		Password: string(hash),
		Nickname: req.Nickname,
	}
	if err := dao.Register(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "数据解析失败"})
		return
	}
	u, err := dao.FindUserByusername(req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "用户名或密码错误"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "用户名或密码错误"})
		return
	}
	token, err := utils.GenToken(u.ID, u.Username, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "登录成功", "token": token})
}

func UserInfo(c *gin.Context) {
	username := c.GetString(middleware.CtxUsername)
	u, err := dao.FindUserByusername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "获取成功", "data": u})
}

func RegisterUserRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", Register)
	rg.POST("/login", Login)
}

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	rg.GET("/me", UserInfo)
}
