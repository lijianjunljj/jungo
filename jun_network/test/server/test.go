package main

import (
	"github.com/lijianjunljj/jungo/jun_boot"
	gate "github.com/lijianjunljj/jungo/jun_gate"
	"github.com/lijianjunljj/jungo/jun_log"
	"github.com/lijianjunljj/jungo/jun_network/boot"
	"github.com/lijianjunljj/jungo/jun_network/conf/ws_server"
	"github.com/lijianjunljj/jungo/jun_network/json"
)

type UserLoginTos struct {
	Token string
}
type UserLoginToc struct {
	Name string
}
type UserLogin struct {
	S *UserLoginTos
	C *UserLoginToc
}
type UserLoginTc struct {
}

func main() {

	jun_boot.Run(func() {
		serverName := "test_gate_server"
		handler := boot.NewHandler(serverName, json.NewProcessor())
		{
			handler.Register(&UserLogin{}, func(state interface{}, args1 ...interface{}) {
				arg := args1[0].([]interface{})
				msg := arg[0]
				agent := arg[1].(gate.Agent)
				jun_log.Debug("msg:%v", msg)
				agent.WriteMsg(&UserLogin{
					C: &UserLoginToc{
						Name: "lijianjun",
					},
				})
			})
		}
		ws_server.ServerConf.WSAddr = "0.0.0.0"
		ws_server.ServerConf.WSPort = "4433"
		boot.ServerStart(serverName, &ws_server.ServerConf, handler)
	})

}
