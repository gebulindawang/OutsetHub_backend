package dao

import (
	"backend/model"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	if err := godotenv.Load(); err != nil{
		panic(err)
	}
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func Register(u *model.User) error {
	result := DB.Create(u)
	return result.Error
}

func FindUserByusername(username string) (model.User,error){
	var user model.User
	err := DB.Where("username = ?",username).First(&user).Error
	return user,err
}


