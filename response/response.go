package response

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Body 所有接口的统一响应外壳
//
//	成功:{"code":0,"msg":"ok","data":...}
//	失败:{"code":400,"msg":"参数错误"}
//
// HTTP 状态码保持语义不变(400/401/403/404/500),
// code 字段与状态码同值,仅作为冗余标记,方便前端统一读取。
type Body struct {
	Code int    `json:"code" example:"0"`
	Msg  string `json:"msg" example:"ok"`
	Data any    `json:"data,omitempty"`
}

// Ok 成功并返回数据
func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: 0, Msg: "ok", Data: data})
}

// OkMsg 成功但无需返回数据
func OkMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Body{Code: 0, Msg: msg})
}

// Fail 业务失败,msg 是给用户看的文案
func Fail(c *gin.Context, status int, msg string) {
	c.JSON(status, Body{Code: status, Msg: msg})
}

// FailErr 业务失败且需要记录底层错误。err 只写入服务端日志,不进入响应体,
// 避免把表名、SQL 错误码等内部细节暴露给客户端。
func FailErr(c *gin.Context, status int, msg string, err error) {
	if err != nil {
		log.Printf("[error] %s %s -> %s: %v", c.Request.Method, c.Request.URL.Path, msg, err)
	}
	c.JSON(status, Body{Code: status, Msg: msg})
}

// Abort 用于中间件:写入响应并终止后续 handler
func Abort(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, Body{Code: status, Msg: msg})
}
