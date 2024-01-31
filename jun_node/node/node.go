package node

import (
	gate "github.com/lijianjunljj/jungo/jun_gate"
)

type Node struct {
	Name  string
	Agent gate.Agent
}

func (n *Node) SendMsg(msg interface{}) {
	if n.Agent != nil {
		n.Agent.WriteMsg(msg)
	}
}
