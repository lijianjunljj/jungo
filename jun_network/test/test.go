package main

import (
	"fmt"
	network "github.com/lijianjunljj/jungo/jun_network"
	"github.com/lijianjunljj/jungo/jun_network/json"
	"github.com/lijianjunljj/jungo/jun_network/test/msg"
)

func main() {
	wsClient := &network.WSClient{
		Addr:     "ws://127.0.0.1:4477",
		NewAgent: NewAgent,
	}
	wsClient.Start()
	wsClient.Wait()
}

type WsAgent struct {
	conn      *network.WSConn
	Processor *json.Processor
}

func (wsa *WsAgent) Run() {
	fmt.Println("run.......")
	wsa.Processor = msg.Processor
	for {
		data, err := wsa.conn.ReadMsg()
		if err != nil {
			fmt.Println("ReadMsg Error:", err)
			return
		}
		msg, err := wsa.Processor.Unmarshal(data)
		if err != nil {
			fmt.Println("Unmarshal Error:", err)
			return
		}
		fmt.Println("ReadMsg Ok:", msg)
	}

	//jun_server.Run()

}
func (wsa *WsAgent) OnClose() {

}
func NewAgent(conn *network.WSConn) network.Agent {
	return &WsAgent{conn: conn}
}
