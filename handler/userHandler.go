package handler

import (
	"backend/dao"
	"backend/model"
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

func RegisterUserRoutes(rg *gin.RouterGroup){
	rg.POST("/register",Register)
}