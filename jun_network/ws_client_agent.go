package network

import (
	"fmt"
	"github.com/lijianjunljj/jungo/jun_log"
	"github.com/lijianjunljj/jungo/jun_network/json"
	"github.com/lijianjunljj/jungo/jun_server"
)

type WsClientAgent struct {
	Conn          *WSConn
	Processor     *json.Processor
	MsgRouterName string
}

func (wsa *WsClientAgent) SendMsg(msg interface{}) {
	wsa.SendMsg(msg)
}
func (wsa *WsClientAgent) Run() {
	for {
		data, err := wsa.Conn.ReadMsg()
		if err != nil {
			fmt.Println("ReadMsg Error:", err)
			return
		}
		msg, err := wsa.Processor.Unmarshal(data)
		if err != nil {
			fmt.Println("Unmarshal Error:", err)
			return
		}
		err = wsa.Processor.Route(msg, wsa)
		if err != nil {
			jun_log.Debug("route message error: %v", err)
			break
		}
		fmt.Println("ReadMsg Ok:", msg)
	}

	//jun_server.Run()

}
func (wsa *WsClientAgent) OnClose() {
	jun_server.Cast(wsa.MsgRouterName, "close", nil)
}
func NewWsClientAgent(msgRouterName string, conn *WSConn, processor *json.Processor) Agent {
	return &WsClientAgent{Conn: conn, Processor: processor, MsgRouterName: msgRouterName}
}
