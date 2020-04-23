package ginn

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const CodeSuccess = 0
const CodeFailed = -1

func RespSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Resp{Code: CodeSuccess, Msg: "success", Data: data})
}

func RespFail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusBadRequest, Resp{Code: code, Msg: msg, Data: nil})
}
