package node

import (
	"errors"
	"fmt"
	gate "github.com/lijianjunljj/jungo/jun_gate"
	"sync"
)

var Nodes sync.Map

func AddNode(nd *Node) error {
	_, ok := Nodes.Load(nd.Name)
	if ok {
		return errors.New("node already exists")
	}
	Nodes.Store(nd.Name, nd)
	return nil
}

func GetNode(nodeName string) (nd *Node) {
	Nodes.Range(func(key, value any) bool {
		fmt.Println("value.(*Node).Name:", value.(*Node).Name)
		fmt.Println("nodeName:", nodeName)
		if value.(*Node).Name == nodeName {
			nd = value.(*Node)
			return false
		}
		return true
	})
	return
}

type Node struct {
	Name  string
	Agent gate.Agent
}

func (n *Node) SendMsg(msg interface{}) {
	if n.Agent != nil {
		n.Agent.WriteMsg(msg)
	}
}
