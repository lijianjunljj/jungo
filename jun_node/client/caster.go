package client

import (
	"fmt"
	"github.com/lijianjunljj/gocommon/misc"
	"github.com/lijianjunljj/jungo/jun_node/msg"
)

type Caster struct {
	client  *Client
	castMsg *msg.Cast
}

func NewCaster(nodeName, distName, key string, data interface{}) *Caster {
	sessionId := misc.GetUnixIDStr()
	return &Caster{
		castMsg: msg.NewCastMsg(sessionId, nodeName, distName, key, data),
		client:  NodeClient}
}

func (that *Caster) Cast() {

	if that.client != nil {
		that.client.WsClientAgent.SendMsg(that.castMsg)
	} else {
		fmt.Println("caster start client is nil")
	}

}
