package server

import (
	"errors"
	gate "github.com/lijianjunljj/jungo/jun_gate"
	"github.com/lijianjunljj/jungo/jun_node/conf"
	"github.com/lijianjunljj/jungo/jun_node/node"
	"github.com/lijianjunljj/jungo/jun_server"
	"reflect"
	"strconv"
	"sync"
)

const (
	ServerName = "`node_server`"
)

type State struct {
	Nodes sync.Map
}

func (s *State) AddNode(nd *node.Node) error {
	_, ok := s.Nodes.Load(nd.Name)
	if ok {
		return errors.New("node already exists")
	}
	s.Nodes.Store(nd.Name, nd)
	return nil
}

type Server struct {
	jun_server.Server
	gate     *gate.Gate
	closeSig chan bool
}

func newServer() jun_server.ModuleBehavior {
	return &Server{Server: jun_server.Server{
		ServerName: ServerName,
		State:      &State{},
	},
	}
}

func Start() {
	//caller.Start(newServer)
}

func (that *Server) Start(interface{}) {
	that.gate = &gate.Gate{
		MaxConnNum:      conf.ServerMaxConnNum,
		PendingWriteNum: conf.PendingWriteNum,
		MaxMsgLen:       conf.MaxMsgLen,
		WSAddr:          conf.ServerListenAddr + ":" + strconv.Itoa(conf.ServerPort),
		HTTPTimeout:     conf.HTTPTimeout,
		CertFile:        conf.CertFile,
		KeyFile:         conf.KeyFile,
		LenMsgLen:       conf.LenMsgLen,
		LittleEndian:    conf.LittleEndian,
		Processor:       Processor,
		AgentServerName: that.ServerName,
	}
	that.closeSig = make(chan bool)
	go that.gate.Run(that.closeSig)
}

func (that *Server) RegisterEvent() {
	that.RegisterCast("NewAgent", rpcNewAgent)
	that.RegisterCast("CloseAgent", rpcCloseAgent)
}

func (that *Server) Terminate(interface{}) {
	that.closeSig <- true
}
func (that *Server) handler(m interface{}, h interface{}) {
	Processor.SetRouter(m, that.ServerName)
	that.RegisterCast(reflect.TypeOf(m), h.(func(interface{}, ...interface{})))
}
