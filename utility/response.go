package utility

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
{
	status:,//5位数状态码，正常下返回10000
	info:,//表示返回的信息，正常下success
	data:,//一个JSON数据，包含要返回的信息，注意，只有正确响应时才返回数据
}
*/

type Response struct {
	Status int         `json:"status"`
	Info   interface{} `json:"info"`           //为什么不是string？因为如果在与数据库字段匹配上，发生了多个错误，返回的err就可能是map之类的
	Data   interface{} `json:"data,omitempty"` //因为并不知道数据的类型是什么
}

func ResponseErr(c *gin.Context, status int) {
	var res = Response{status, GetInfo(status), nil}
	c.JSON(http.StatusOK, &res) //可以传递指针，优化性能
}

// 这个data interface{}必须是结构体或者map吗，毕竟序列化必须是一对一的嘛
func ResponseSuccess(c *gin.Context, data interface{}) {
	var res = Response{StatusSuccess, GetInfo(StatusSuccess), data}
	c.JSON(http.StatusOK, &res)
}
