package client

import (
	"github.com/lijianjunljj/gocommon/misc"
	"github.com/lijianjunljj/jungo/jun_log"
	"github.com/lijianjunljj/jungo/jun_node/msg"
	"github.com/lijianjunljj/jungo/jun_server"
)

type Caller struct {
	jun_server.Server
	client  *Client
	callMsg *msg.Call
}

func NewCaller(nodeName, distName, key string, data interface{}) *Caller {
	sessionId := misc.GetUnixIDStr()
	return &Caller{
		Server: jun_server.Server{
			ServerName: sessionId,
			CloseSig:   make(chan jun_server.ExitSig),
			State:      sessionId,
		},
		callMsg: msg.NewCallMsg(sessionId, nodeName, distName, key, data),
		client:  NodeClient}
}

func (that *Caller) Start(interface{}) {
	that.client.WsClientAgent.SendMsg(that.callMsg)

}
func (that *Caller) RegisterEvent() {
	that.RegisterCast("reply", func(iState interface{}, args ...interface{}) {
		data := args[0]
		that.Exit(jun_server.ExitReasonNormal, data)
	})
}

func (that *Caller) Terminate(iState interface{}) {
	jun_log.Debug("client caller terminate caller:", iState.(string))
}
