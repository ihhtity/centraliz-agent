package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 带自定义消息的成功响应
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, msg string, data ...interface{}) {
	if len(data) == 0 {
		data = nil
	}

	// 根据业务错误码返回对应的HTTP状态码
	var httpStatus int
	switch {
	case code >= 500:
		httpStatus = http.StatusInternalServerError
	case code >= 400:
		httpStatus = http.StatusBadRequest
	default:
		httpStatus = http.StatusOK
	}

	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// Error 服务器错误响应
func Error(c *gin.Context, msg string, data ...interface{}) {
	if len(data) == 0 {
		data = nil
	}

	c.JSON(http.StatusInternalServerError, Response{
		Code: 500,
		Msg:  msg,
		Data: data,
	})
}
