package main

import (
	"fmt"
	"github.com/lijianjunljj/jungo/jun_boot"
	"github.com/lijianjunljj/jungo/jun_log"
	"github.com/lijianjunljj/jungo/jun_node"
	"github.com/lijianjunljj/jungo/jun_node/conf"
	"github.com/lijianjunljj/jungo/jun_server"
	"time"
)

func main() {

	jun_boot.Run(func() {
		jun_node.Start()
		Start()
		if conf.NodeName == "node2" {
			time.Sleep(10 * time.Second)
			ret := jun_node.Call("node1", "TestServer", "test_call", "hello world!")
			jun_log.Debug("ret:%v", ret.Replay)
			jun_node.Cast("node1", "TestServer", "test_cast", "cast message!")
		}

	})

}

type TestServer struct {
	jun_server.Server
}

func newServer() jun_server.ModuleBehavior {
	return &TestServer{Server: jun_server.Server{
		ServerName: "TestServer",
	},
	}
}

func Start() {
	jun_server.Start(newServer)
}
func (that *TestServer) Start(interface{}) {

}
func (that *TestServer) RegisterEvent() {
	that.RegisterCall("test_call", func(iState interface{}, args ...interface{}) *jun_server.CallRet {
		return &jun_server.CallRet{Replay: args[0]}
	})
	that.RegisterCast("test_cast", func(iState interface{}, args ...interface{}) {
		fmt.Println("test_cast args:", args)
	})
}
func (that *TestServer) Terminate(interface{}) {
}
