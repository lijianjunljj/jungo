package main

import (
	"github.com/lijianjunljj/jungo/jun_boot"
	"github.com/lijianjunljj/jungo/jun_network/boot"
	"github.com/lijianjunljj/jungo/jun_network/conf/ws_server"
	"github.com/lijianjunljj/jungo/jun_network/json"
)

func main() {

	jun_boot.Run(func() {
		serverName := "test_gate_server"
		handler := boot.NewHandler(serverName, json.NewProcessor())
		ws_server.ServerConf.WSAddr = "0.0.0.0"
		ws_server.ServerConf.WSPort = "4433"
		boot.ServerStart(serverName, &ws_server.ServerConf, handler)
	})

}
