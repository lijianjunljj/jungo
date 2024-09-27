package client

import (
	"github.com/lijianjunljj/gocommon/misc"
	"github.com/lijianjunljj/jungo/jun_log"
	"github.com/lijianjunljj/jungo/jun_node/msg"
	"github.com/lijianjunljj/jungo/jun_server"
)

type Caster struct {
	jun_server.Server
	client  *Client
	castMsg *msg.Cast
}

func NewCaster(nodeName, distName, key string, data interface{}) *Caster {
	sessionId := misc.GetUnixIDStr()
	return &Caster{
		Server: jun_server.Server{
			ServerName: sessionId,
			CloseSig:   make(chan jun_server.ExitSig),
			State:      sessionId,
		},
		castMsg: msg.NewCastMsg(sessionId, nodeName, distName, key, data),
		client:  NodeClient}
}

func (that *Caster) Start(interface{}) {
	that.client.WsClientAgent.SendMsg(that.castMsg)

}
func (that *Caster) RegisterEvent() {
	that.RegisterCast("reply", func(iState interface{}, args ...interface{}) {
		data := args[0]
		that.Exit(jun_server.ExitReasonNormal, data)
	})
}

func (that *Caster) Terminate(iState interface{}) {
	jun_log.Debug("client caller terminate caller:", iState.(string))
}
