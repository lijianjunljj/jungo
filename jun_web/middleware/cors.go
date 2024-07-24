package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization,X_Requested_With,x-requested-with,Signature,AppId")
		ctx.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Type,Authorization,X_Requested_With,Signature,AppId")

		if ctx.Request.Method == http.MethodOptions {
			ctx.JSON(200, nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
