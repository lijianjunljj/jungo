package msg

import (
	"fmt"
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
	fmt.Println("that.DistModName, that.Key, that.Msg:", that.DistModName, that.Key, that.Msg)
	callRet := jun_server.Call(that.DistModName, that.Key, that.Msg)
	fmt.Println("callRet:", callRet)
	that.SetTransportToBack()
	that.Reply = callRet
}
