package main

import (
	"backend/dao"
	"backend/router"
)

func main() {
	if err := dao.InitDB(); err != nil{
		panic(err)
	}
	r := router.SetupRouter()
	r.Run("127.0.0.1:8080")
}