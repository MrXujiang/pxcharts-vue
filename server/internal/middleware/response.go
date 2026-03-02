package middleware

import (
	"mvtable/internal/pkg/errorx"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
// @Description 统一响应结构
type Response struct {
	Data  interface{} `json:"data"`  // 响应数据
	Msg   string      `json:"msg"`   // 响应消息
	State int         `json:"state"` // 响应状态码
}

func ResJSON(c *gin.Context, code int, data any, msg string, state int) {
	c.JSON(code, gin.H{
		"data":  data,
		"msg":   msg,
		"state": state,
	})
	c.Abort()
}

func ResOK(c *gin.Context, data any, msg string) {
	ResJSON(c, http.StatusOK, data, msg, http.StatusOK)
}

func ResErr(c *gin.Context, err error) {
	var reqErr *errorx.Error

	if e, ok := errorx.As(err); ok && e.Status <= 100000 {
		reqErr = e
	} else if ok && e.Status > 100000 {
		ResJSON(c, http.StatusOK, nil, e.Message, e.Status)
		return
	} else {
		reqErr = errorx.InternalServerError("服务内部错误")
	}

	ResJSON(c, reqErr.Status, nil, reqErr.Message, reqErr.Status)
}
