package msg

import (
	"github.com/lijianjunljj/jungo/jun_network/json"
	"github.com/lijianjunljj/jungo/jun_node/conf"
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
	Processor.Register(&Cast{})
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

type IMsg interface {
	IsTargetNode() bool
	GetTransportNode() string
	SetTransportToBack()
	GetTransportType() int
	GetDistNodeName() string
	GetSrcNodeName() string
	TransportToBack()
}

type Msg struct {
	TransportType int
	Key           string
	SessionId     string
	SrcNodeName   string
	DistNodeName  string
	DistModName   string
	Msg           interface{}
}

func (that *Msg) GetTransportType() int {
	return that.TransportType
}
func (that *Msg) GetDistNodeName() string {
	return that.DistNodeName
}
func (that *Msg) GetSrcNodeName() string {
	return that.SrcNodeName
}

func (that *Msg) IsTargetNode() bool {
	nodeName := that.GetTransportNode()
	return nodeName == conf.NodeName
}

func (that *Msg) GetTransportNode() string {
	switch that.TransportType {
	case CallTransportTypeEnter:
		return that.DistNodeName
	case CallTransportTypeBack:
		return that.SrcNodeName
	}
	return ""
}
func (that *Msg) SetTransportToBack() {
	that.TransportType = CallTransportTypeBack
}
