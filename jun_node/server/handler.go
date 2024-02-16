package server

import (
	"fmt"
	gate "github.com/lijianjunljj/jungo/jun_gate"
	"github.com/lijianjunljj/jungo/jun_log"
	"github.com/lijianjunljj/jungo/jun_node/client"
	"github.com/lijianjunljj/jungo/jun_node/conf"
	"github.com/lijianjunljj/jungo/jun_node/msg"
	"github.com/lijianjunljj/jungo/jun_node/node"
)

func (that *Server) InitHandler() {
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	that.handler(&msg.LoginTos{}, that.HandleLogin)
	that.handler(&msg.Call{}, that.HandleCall)
}
func (that *Server) HandleLogin(iState interface{}, args ...interface{}) {
	jun_log.Debug("HandleLogin....", args[0])
	arg := args[0].([]interface{})
	a := arg[1].(gate.Agent)
	m := arg[0].(*msg.LoginTos)
	fmt.Println(a, m)
	state := iState.(*State)
	nd := &node.Node{Name: m.NodeName, Agent: a}
	a.SetUserData(nd)
	err := state.AddNode(nd)
	if err != nil {
		nd.SendMsg(&msg.MsgComToc{Code: 1, Msg: err.Error()})
		return
	}
	nd.SendMsg(&msg.MsgComToc{})
}

func (that *Server) HandleCall(iState interface{}, args ...interface{}) {
	jun_log.Debug("HandleCall....")
	arg := args[0].([]interface{})
	a := arg[1].(gate.Agent)
	m := arg[0].(*msg.Call)
	fmt.Println(a, m)

	iNode := a.UserData()
	_ = iNode.(*node.Node)

	state := iState.(*State)
	var distNode *node.Node
	switch m.TransportType {
	case msg.CallTransportTypeEnter:
		if m.DistNodeName == conf.NodeName {
			m.TransportToBack()
			a.WriteMsg(m)
			return
		}
		distNode = state.GetNode(m.DistNodeName)
		break
	case msg.CallTransportTypeBack:
		distNode = state.GetNode(m.SrcNodeName)
		break
	}
	if distNode != nil {
		distNode.SendMsg(m)
	} else {
		if client.NodeClient != nil {
			client.NodeClient.WsClientAgent.SendMsg(m)
		}
	}

}
func rpcNewAgent(state interface{}, args ...interface{}) {
	fmt.Println("args:", args)
	a := args[0].(gate.Agent)
	fmt.Println("NewAgent............")
	_ = a
}

func rpcCloseAgent(state interface{}, args ...interface{}) {
	fmt.Println("CloseAgent............")
	a := args[0].(gate.Agent)
	fmt.Println(a)
	userData := a.UserData()
	fmt.Println(userData)

}
