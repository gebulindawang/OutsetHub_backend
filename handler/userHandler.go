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

func Register(c *gin.Context){
	var user model.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	newPassword,err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"err":err.Error()})
		return
	}
	user.Password = string(newPassword)
	if err:= dao.Register(&user); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"err":err.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{"msg":"注册成功"})
}

func Login(c *gin.Context){
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"err":"数据解析失败"})
		return
	}
	u,err := dao.FindUserByusername(user.Username)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"err":err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password),[]byte(user.Password));err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"err":err.Error()})
		return
	}
	token ,err:= utils.GenToken(u.ID,u.Username,u.Role)
	if err != nil{
		panic(err)	
	}

	c.JSON(http.StatusOK,gin.H{"msg":"登录成功","token":token})
}

func UserInfo(c *gin.Context){
	username := c.GetString(middleware.CtxUsername)
	u,err := dao.FindUserByusername(username)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK,gin.H{"data":u})

	c.JSON(200,gin.H{"msg":"获取成功"})
}


func RegisterUserRoutes(rg *gin.RouterGroup){
	rg.POST("/register",Register)
	rg.POST("/login",Login)
}

func RegisterAuthRoutes(rg *gin.RouterGroup){
	rg.GET("/me",UserInfo)
}