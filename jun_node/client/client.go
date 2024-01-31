package client

import (
	"github.com/lijianjunljj/jungo/jun_server"
)

const (
	ServerName = "`node_client`"
)

type State struct {
}

type Client struct {
	jun_server.Server
}

func newServer() jun_server.ModuleBehavior {
	return &Client{Server: jun_server.Server{
		ServerName: ServerName,
		State:      &State{},
	},
	}
}

func Start() {
	//caller.Start(newServer)
}

func (that *Client) Start(interface{}) {

}

func (that *Client) RegisterEvent() {
}

func (that *Client) Terminate(interface{}) {

}
