package msg

import (
	"github.com/lijianjunljj/jungo/jun_network/json"
	"github.com/lijianjunljj/jungo/jun_node/conf"
	"github.com/lijianjunljj/jungo/jun_server"
)

const (
	CallTransportTypeEnter = 1
	CallTransportTypeBack  = 2
)

func NewProcessor() *json.Processor {
	var Processor = json.NewProcessor()
	Processor.Register(&LoginTos{})
	Processor.Register(&LoginToc{})
	Processor.Register(&Call{})
	Processor.Register(&MsgComToc{})

	return Processor
}

type MsgComToc struct {
	Code int8   `json:"code"`
	Msg  string `json:"msg"`
}
type LoginTos struct {
	Cookie   string
	NodeName string
}

type LoginToc struct {
	MsgComToc
}

func NewCallMsg(sessionId, nodeName, distName, key string, data interface{}) *Call {
	return &Call{SessionId: sessionId,
		Key:           key,
		SrcNodeName:   conf.NodeName,
		DistNodeName:  nodeName,
		DistModName:   distName,
		TransportType: CallTransportTypeEnter,
		Msg:           data}
}

type Call struct {
	TransportType int8
	Key           string
	SessionId     string
	SrcNodeName   string
	DistNodeName  string
	DistModName   string
	Msg           interface{}
	Reply         interface{}
}

func (that *Call) IsTargetNode() bool {
	nodeName := that.GetTransportNode()
	return nodeName == conf.NodeName
}

func (that *Call) GetTransportNode() string {
	switch that.TransportType {
	case CallTransportTypeEnter:
		return that.DistNodeName
	case CallTransportTypeBack:
		return that.SrcNodeName
	}
	return ""
}
func (that *Call) TransportToBack() {
	callRet := jun_server.Call(that.DistModName, that.Key, that.Msg)
	that.TransportType = CallTransportTypeBack
	that.Reply = callRet
}
