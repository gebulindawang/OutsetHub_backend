package main

import (
	"backend/dao"
	"backend/router"
)

// 通用接口信息,由 swag init 读取并生成 docs/
//
//	@title						起点社区 OutsetHub API
//	@version					1.0
//	@description				起点社区后端接口文档。除注册、登录外,所有接口都需要在请求头携带 Bearer token。
//	@description				使用方式:先调用 POST /api/user/login 拿到 token,再点击右上角 Authorize 填入 token 原文(不需要自己加 Bearer 前缀)。
//	@description				统一响应约定:成功为 {msg, data} 或 {msg};业务错误为 {msg, err};中间件错误为 {code, msg}。
//	@host						127.0.0.1:8080
//	@BasePath					/api
//	@schemes					http
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				格式:Bearer <token>
func main() {
	if err := dao.InitDB(); err != nil {
		panic(err)
	}
	r := router.SetupRouter()
	r.Run("127.0.0.1:8080")
}
