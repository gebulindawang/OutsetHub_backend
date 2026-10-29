package dao

import (
	"backend/model"

	"gorm.io/gorm"
)
var DB *gorm.DB
	
func Register(u *model.User) error {
	result := DB.Create(u)
	return result.Error
}

func FindUserByusername(username string) (model.User,error){
	var user model.User
	err := DB.Where("username = ?",username).First(&user).Error
	return user,err
}




