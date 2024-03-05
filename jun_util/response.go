package jun_util

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	//CodeSuccess 请求成功响应码
	CodeSuccess = "0000"
	//CodeFail 请求失败响应码
	CodeFail         = "1001"
	CodeTokenExpired = "1002"
	//CodeException 系统异常响应码
	CodeException = "1003"
	CodeUserBlack = "1004"
)

// Response 响应结构体
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// Success 响应成功
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, Response{
		Code:    CodeSuccess,
		Message: "成功",
		Data:    data,
	})
}

// Fail 响应错误
func Fail(ctx *gin.Context, err error, codes ...string) {
	code := CodeFail
	msg := err.Error()
	//fmt.Println(err, msg)
	if len(codes) > 0 {
		code = codes[0]
	}
	if _, isValidationErrors := err.(validator.ValidationErrors); isValidationErrors {
		msg = err.Error() // "参数校验不通过"
	}
	if _, isJSONError := err.(*json.UnmarshalTypeError); isJSONError {
		msg = err.Error() // "字段类型不匹配"
	}
	if strings.Contains(msg, "rpc error") {
		errs := strings.Split(msg, "=")
		msg = errs[len(errs)-1]
	}
	ctx.JSON(200, Response{
		Code:    code,
		Message: msg,
		Data:    err.Error(),
	})
}
