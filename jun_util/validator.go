package jun_util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

var validate *validator.Validate
var validateTrans ut.Translator

// InitValid 初始化校验器
func InitValid() {
	zh_ch := zh.New()
	validate = validator.New()
	uni := ut.New(zh_ch)
	validateTrans, _ := uni.GetTranslator("zh")
	//验证器注册翻译器
	zh_translations.RegisterDefaultTranslations(validate, validateTrans)
	// 注册自定义验证
	validate.RegisterValidation("nameValid", nameValid)
}

// Validate 参数校验
func Validate(s interface{}, ctx *gin.Context) error {
	if validate == nil {
		InitValid()
	}
	_, err := Bind(ctx, s)
	if err != nil {
		return err
	}
	if err = validate.Struct(s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(validateTrans))
		}
		return err
	}
	return nil
}

// 自定义校验函数
func nameValid(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if val == "admin" {
		return false
	}
	return true
}
