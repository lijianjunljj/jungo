package jun_node

import (
	"fmt"
	"github.com/lijianjunljj/jungo/jun_log"
	"github.com/lijianjunljj/jungo/jun_node/client"
	"github.com/lijianjunljj/jungo/jun_node/conf"
	"github.com/lijianjunljj/jungo/jun_node/server"
	"github.com/lijianjunljj/jungo/jun_server"
	"github.com/lijianjunljj/jungo/jun_util"
)

const (
	ServerName = "`node_scheduler`"
)

type State struct {
}

func Start() {
	jun_server.Start(newServer)
}

type Scheduler struct {
	jun_server.Server
	client *client.Client
	server *server.Server
}

func newServer() jun_server.ModuleBehavior {
	return &Scheduler{Server: jun_server.Server{
		ServerName: ServerName,
		State:      &State{},
	},
	}
}

func (that *Scheduler) Start(interface{}) {
	port := jun_util.CheckPorts([]int{conf.ServerPort})

	fmt.Println("port:", port)
	if port != "" {
		if conf.CenterNodeHost == "" {
			jun_log.Debug("开始启动中心节点")
		} else {
			jun_log.Debug("开始启动主节点")
		}
		server.Start()
		if conf.CenterNodeHost != "" {
			client.Start(conf.CenterNodeHost)
		}
	} else {
		jun_log.Debug("开始启动从节点")
		client.Start("127.0.0.1")
	}
}

func (that *Scheduler) RegisterEvent() {
}

func (that *Scheduler) Terminate(interface{}) {

}
