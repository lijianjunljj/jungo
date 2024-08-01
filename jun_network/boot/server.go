package boot

import (
	"fmt"
	gate "github.com/lijianjunljj/jungo/jun_gate"
	"github.com/lijianjunljj/jungo/jun_network/conf/ws_server"
	"github.com/lijianjunljj/jungo/jun_server"
)

type GateServer struct {
	jun_server.Server
	gate     *gate.Gate
	closeSig chan bool
	Conf     *ws_server.Config
	handler  Handler
}

func newServer(serverName string, conf *ws_server.Config, handler Handler) jun_server.ModuleBehavior {
	return &GateServer{Server: jun_server.Server{ServerName: serverName}, Conf: conf, handler: handler}
}

func ServerStart(serverName string, conf *ws_server.Config, handler Handler) {
	jun_server.Start(func() jun_server.ModuleBehavior {
		return newServer(serverName, conf, handler)
	})
}

func (that *GateServer) Start(interface{}) {
	that.gate = &gate.Gate{
		MaxConnNum:      that.Conf.MaxConnNum,
		PendingWriteNum: that.Conf.PendingWriteNum,
		MaxMsgLen:       that.Conf.MaxMsgLen,
		WSAddr:          that.Conf.WSAddr + ":" + that.Conf.WSPort,
		HTTPTimeout:     that.Conf.HTTPTimeout,
		CertFile:        that.Conf.CertFile,
		KeyFile:         that.Conf.KeyFile,
		TCPAddr:         that.Conf.TCPAddr,
		LenMsgLen:       that.Conf.LenMsgLen,
		LittleEndian:    that.Conf.LittleEndian,
		Processor:       that.handler.Processor,
		AgentServerName: that.ServerName,
	}
	that.closeSig = make(chan bool)
	go that.gate.Run(that.closeSig)
}
func (that *GateServer) RegisterEvent() {
	that.RegisterCast("NewAgent", rpcNewAgent)
	that.RegisterCast("CloseAgent", rpcCloseAgent)
	that.InitHandler()
}
func (that *GateServer) Terminate(interface{}) {
	fmt.Println("GateServer terminate......")
	that.closeSig <- true
}

func (that *GateServer) InitHandler() {
	that.handler.RouterAll(that.Server)
}
func rpcNewAgent(state interface{}, args ...interface{}) {
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
	if userData != nil {
	}
}
