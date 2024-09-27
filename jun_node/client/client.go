package client

import (
	"fmt"
	"github.com/lijianjunljj/jungo/jun_log"
	network "github.com/lijianjunljj/jungo/jun_network"
	"github.com/lijianjunljj/jungo/jun_network/json"
	"github.com/lijianjunljj/jungo/jun_node/conf"
	"github.com/lijianjunljj/jungo/jun_node/msg"
	"github.com/lijianjunljj/jungo/jun_node/node"
	"github.com/lijianjunljj/jungo/jun_server"
	"reflect"
	"runtime/debug"
	"strconv"
)

const (
	ServerName = "`node_client`"
)

var NodeClient *Client

type State struct {
}

type Client struct {
	jun_server.Server
	Processor     *json.Processor
	WsClientAgent *network.WsClientAgent
	masterHost    string
}

func Start(masterHost string) {
	newServerFuc := func() jun_server.ModuleBehavior {
		NodeClient = &Client{Server: jun_server.Server{
			ServerName: ServerName,
			State:      &State{},
		},
			Processor:  msg.NewProcessor(),
			masterHost: masterHost,
		}
		return NodeClient
	}
	jun_server.Start(newServerFuc)
}

func (that *Client) Start(interface{}) {
	that.Connect()
}

func (that *Client) Connect() {
	newAgent := func(conn *network.WSConn) network.Agent {
		that.WsClientAgent = &network.WsClientAgent{Conn: conn, MsgRouterName: ServerName, Processor: that.Processor}
		return that.WsClientAgent
	}
	wsClient := &network.WSClient{
		Addr:     "ws://" + that.masterHost + ":" + strconv.Itoa(conf.ServerPort),
		NewAgent: newAgent,
	}
	wsClient.Start()
}

func (that *Client) RegisterEvent() {
	that.RegisterCast("close", func(iState interface{}, args ...interface{}) {
		jun_log.Debug("从节点断开，开始重连...")
		that.Connect()
	})

	that.RegisterCast("NewAgent", func(iState interface{}, args ...interface{}) {
		//state := iState.(*State)
		jun_log.Debug("从节点连接成功，开始登录主节点...")
		wsa := args[0].(*network.WsClientAgent)
		fmt.Println("iState:", iState)
		cookie := conf.Cookie
		err := wsa.SendMsg(&msg.LoginTos{Cookie: cookie, NodeName: conf.NodeName})
		if err != nil {
			fmt.Println("err:", err)
			return
		}
	})

	that.handlerMsg(&msg.Call{}, func(iState interface{}, iMsg interface{}, wsa *network.WsClientAgent) {
		m := iMsg.(*msg.Call)
		switch m.TransportType {
		case msg.CallTransportTypeBack:
			if m.SrcNodeName == conf.NodeName {
				callRet := jun_server.Call(m.DistModName, m.Key, m.Msg)
				jun_server.Cast(m.SessionId, "reply", callRet)
				return
			}
			fmt.Println("m.SrcNodeName:", m.SrcNodeName)
			distNode := node.GetNode(m.SrcNodeName)
			fmt.Println("distNode:", distNode)
			if distNode != nil {
				distNode.SendMsg(m)
			}
			break
			//case msg.CallTransportTypeEnter:
			//	wsa.SendMsg(m)
			//	break
		}

	})
}

func (that *Client) Terminate(interface{}) {

}

func (that *Client) handlerMsg(m interface{}, fn func(interface{}, interface{}, *network.WsClientAgent)) {
	that.Processor.SetRouter(m, that.ServerName)

	f := func(state interface{}, args1 ...interface{}) {
		go func() {
			defer func() {
				if msg := recover(); msg != nil {
					jun_log.Error("panic信息:%v", msg)
					jun_log.Error("调用栈:%v", string(debug.Stack()))
				}
			}()
			arg := args1[0].([]interface{})
			msg := arg[0]
			agent := arg[1].(*network.WsClientAgent)
			fn(state, msg, agent)
		}()
	}
	that.RegisterCast(reflect.TypeOf(m), f)
}
