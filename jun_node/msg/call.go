package msg

import (
	"github.com/lijianjunljj/jungo/jun_node/conf"
	"github.com/lijianjunljj/jungo/jun_server"
)

func NewCallMsg(sessionId, nodeName, distName, key string, data interface{}) *Call {
	return &Call{
		Msg: Msg{SessionId: sessionId,
			Key:           key,
			SrcNodeName:   conf.NodeName,
			DistNodeName:  nodeName,
			DistModName:   distName,
			TransportType: CallTransportTypeEnter,
			Msg:           data,
		},
	}
}

type Call struct {
	Msg
	Reply interface{}
}

func (that *Call) TransportToBack() {
	callRet := jun_server.Call(that.DistModName, that.Key, that.Msg)
	that.SetTransportToBack()
	that.Reply = callRet
}
