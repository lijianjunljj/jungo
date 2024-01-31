package server

import (
	"fmt"
	gate "github.com/lijianjunljj/jungo/jun_gate"
	"github.com/lijianjunljj/jungo/jun_node/node"
)

func (that *Server) Init() {
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	that.handler(&LoginTos{}, HandleLogin)
}
func HandleLogin(iState interface{}, args ...interface{}) {
	a := args[0].(gate.Agent)
	m := args[1].(*LoginTos)
	fmt.Println(a, m)
	state := iState.(*State)
	nd := &node.Node{Name: m.NodeName, Agent: a}
	err := state.AddNode(nd)
	if err != nil {
		nd.SendMsg(&MsgComToc{Code: 1, Msg: err.Error()})
		return
	}
	nd.SendMsg(&MsgComToc{})
}

func rpcNewAgent(state interface{}, args ...interface{}) {
	fmt.Println("args:", args)
	a := args[0].(gate.Agent)
	fmt.Println("NewAgent............")
	fmt.Println(a)
	_ = a
}

func rpcCloseAgent(state interface{}, args ...interface{}) {
	fmt.Println("CloseAgent............")
	a := args[0].(gate.Agent)
	fmt.Println(a)
	userData := a.UserData()
	fmt.Println(userData)

}
