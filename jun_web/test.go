package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lijianjunljj/jungo/jun_boot"
	"github.com/lijianjunljj/jungo/jun_web/jun_web"
)

func Register(app *gin.Engine) {
	router := app.Group("test")
	{
		router.GET("/index", func(ctx *gin.Context) {
			ctx.String(200, "hellow")
		})
	}
}
func main() {

	jun_boot.Run(func() {
		jun_web.Start([]string{}, Register)
	})
}
