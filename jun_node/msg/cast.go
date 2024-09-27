package msg

import (
	"github.com/lijianjunljj/jungo/jun_node/conf"
	"github.com/lijianjunljj/jungo/jun_server"
)

func NewCastMsg(sessionId, nodeName, distName, key string, data interface{}) *Cast {
	return &Cast{
		Msg{SessionId: sessionId,
			Key:           key,
			SrcNodeName:   conf.NodeName,
			DistNodeName:  nodeName,
			DistModName:   distName,
			TransportType: CallTransportTypeEnter,
			Msg:           data,
		},
	}
}

type Cast struct {
	Msg
}

func (that *Cast) TransportToBack() {
	jun_server.Cast(that.DistModName, that.Key, that.Msg)
	that.SetTransportToBack()
}
