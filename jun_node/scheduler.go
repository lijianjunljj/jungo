package jun_node

import (
	"github.com/lijianjunljj/jungo/jun_node/client"
	"github.com/lijianjunljj/jungo/jun_node/conf"
	"github.com/lijianjunljj/jungo/jun_node/server"
	"github.com/lijianjunljj/jungo/jun_server"
	"github.com/lijianjunljj/jungo/jun_util"
)

const (
	ServerName = "`node_monitor`"
)

type State struct {
}

func Start() {
	//caller.Start(newServer)
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

	if port == "" {

	}
}

func (that *Scheduler) RegisterEvent() {
}

func (that *Scheduler) Terminate(interface{}) {

}
