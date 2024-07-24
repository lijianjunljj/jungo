package jun_web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lijianjunljj/jungo/jun_server"
	"github.com/lijianjunljj/jungo/jun_web/conf"
	"github.com/lijianjunljj/jungo/jun_web/middleware"
	"strconv"
)

const (
	WebServerName = "web_server"
)

type WebServerRouterFunc = func(app *gin.Engine)

type WebServer struct {
	jun_server.Server
	openURL []string
	routers []WebServerRouterFunc
}

func newServer(openURL []string, routers []WebServerRouterFunc) jun_server.ModuleBehavior {
	return &WebServer{
		openURL: openURL,
		routers: routers,
		Server:  jun_server.Server{ServerName: WebServerName},
	}
}
func Start(openURL []string, routers ...func(app *gin.Engine)) {
	jun_server.Start(func() jun_server.ModuleBehavior {
		return newServer(openURL, routers)
	})

}
func (that *WebServer) Start(interface{}) {
	fmt.Println("start webserver.......")
	app := gin.Default()
	app.Use(middleware.Cors())
	app.Use(middleware.Jwt(that.openURL))
	for _, v := range that.routers {
		v(app)
	}
	HttpPort := strconv.FormatInt(int64(conf.Server.HttpPort), 10)
	go func() {
		err := app.Run(":" + HttpPort)
		if err != nil {

		}
	}()
}
func (that *WebServer) RegisterEvent() {

}
func (that *WebServer) Terminate(interface{}) {
	fmt.Println("WebServer terminate......")
}
