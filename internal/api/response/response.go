package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

type Data struct {
	Data interface{} `json:"data"`
}

type List[T any] struct {
	List []T `json:"list"`
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "查询成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithData(data interface{}, c *gin.Context) {
	Result(ERROR, data, "操作失败", c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

func FailToParameter(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{
		Code: ERROR,
		Data: map[string]interface{}{},
		Msg:  "参数错误: " + err.Error(),
	})
}

func FailToInternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Code: ERROR,
		Data: map[string]interface{}{},
		Msg:  "内部错误: " + err.Error(),
	})
}
