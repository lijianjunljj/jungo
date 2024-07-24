package middleware

import (
	"errors"
	"github.com/lijianjunljj/jungo/jun_util"
	"github.com/lijianjunljj/jungo/jun_web/conf"
	"github.com/lijianjunljj/jungo/jun_web/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth 认证中间件
func Jwt(openURL []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !conf.Server.EnableJwt {
			ctx.Next()
			return
		}
		token := ctx.Request.Header.Get("Authorization")
		isOpen := false
		for _, item := range openURL {
			if strings.Contains(ctx.Request.URL.Path, item) {
				isOpen = true
				break
			}
		}
		if isOpen {
			ctx.Next()
			return
		}
		if token == "" || token == "xxxx" {
			util.Fail(ctx, errors.New("Token不能为空2"), util.CodeTokenExpired)
			ctx.Abort()
			return
		}
		res, ok := jun_util.JWTDecrypt(token, conf.Server.JwtSecret)
		if !ok {
			util.Fail(ctx, errors.New("Token过期"), util.CodeTokenExpired)
			ctx.Abort()
			return
		}
		ctx.Set("user", res)
		ctx.Set("userID", res["id"])
		ctx.Next()
	}
}
