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
	that.handler(&msg.Cast{}, that.HandleCast)

}
func (that *Server) HandleLogin(iState interface{}, args ...interface{}) {
	jun_log.Debug("HandleLogin....", args[0])
	arg := args[0].([]interface{})
	a := arg[1].(gate.Agent)
	m := arg[0].(*msg.LoginTos)
	fmt.Println("m.NodeName:", m.NodeName)
	nd := &node.Node{Name: m.NodeName, Agent: a}
	a.SetUserData(nd)
	err := node.AddNode(nd)
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
	that.msgProcess(iState, a, m)
}

func (that *Server) msgProcess(iState interface{}, a gate.Agent, m msg.IMsg) {
	iNode := a.UserData()
	_ = iNode.(*node.Node)

	var distNode *node.Node
	transportType := m.GetTransportType()
	distNodeName := m.GetDistNodeName()
	srcNodeName := m.GetSrcNodeName()
	fmt.Println("transportType,srcNodeName,distNodeName:", transportType, srcNodeName, distNodeName)

	switch transportType {
	case msg.CallTransportTypeEnter:
		if distNodeName == conf.NodeName {
			m.TransportToBack()

			fmt.Println("mm:", m)
			a.WriteMsg(m)
			return
		}
		distNode = node.GetNode(distNodeName)
		break
	case msg.CallTransportTypeBack:
		distNode = node.GetNode(srcNodeName)
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

func (that *Server) HandleCast(iState interface{}, args ...interface{}) {
	jun_log.Debug("HandleCast....")
	arg := args[0].([]interface{})
	a := arg[1].(gate.Agent)
	m := arg[0].(*msg.Cast)
	that.msgProcess(iState, a, m)
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
