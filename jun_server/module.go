package jun_server

import (
	"errors"
	"fmt"
	"github.com/lijianjunljj/jungo/jun_log"
	"github.com/lijianjunljj/jungo/jun_timer"
	"runtime/debug"
	"sync"
	"time"
)

type Server struct {
	ModuleBehavior
	ServerName string
	m          *Module
	CloseSig   chan ExitSig
	State      interface{}
}

func (s *Server) Exit(Reason string, data interface{}) {
	Exit(s.ServerName, Reason, data)
}

func (s *Server) Stop() {
	Stop(s.ServerName)
}

func (s *Server) GetServerName() string {
	return s.ServerName
}
func (s *Server) SetServerName(serverName string) {
	s.ServerName = serverName
}

func (s *Server) GetState() interface{} {
	return s.State
}
func (s *Server) SetModule(m *Module) {
	s.m = m
}
func (s *Server) GetCloseSig() chan ExitSig {
	return s.CloseSig
}

func (s *Server) RegisterCast(key interface{}, f func(interface{}, ...interface{})) {
	s.m.RegisterCast(key, f)
}
func (s *Server) RegisterCall(key interface{}, f func(interface{}, ...interface{}) *CallRet) {
	s.m.RegisterCall(key, f)
}

type NewServer func() ModuleBehavior

var mods sync.Map

var wg sync.WaitGroup

func Start(newServer NewServer) {
	serv := newServer()
	name := serv.GetServerName()
	closeSig := serv.GetCloseSig()
	state := serv.GetState()
	Run(name, serv, closeSig, state)
}

func RunServer(serv ModuleBehavior) error {
	name := serv.GetServerName()
	closeSig := serv.GetCloseSig()
	state := serv.GetState()
	err := Run(name, serv, closeSig, state)
	if err != nil {
		return err
	}
	return nil
}

func Run(name string, serv ModuleBehavior, closeSig chan ExitSig, state interface{}) error {

	//closeSig := make(chan ExitSig)
	if name == "" {
		err := errors.New("server name not empty")
		jun_log.Error(err.Error())
		return err
	}
	_, ok := mods.Load(name)
	if ok {
		err := errors.New("module has started")
		jun_log.Error(err.Error())
		return err
	}
	go func() {
		wg.Add(1)
		m := &Module{State: state,
			ChanCall:    make(chan CallInfo),
			ChanCast:    make(chan CastInfo, 100),
			ChanCallRet: make(chan CallRet),
			ChanExit:    make(chan ExitSig),
			CastRouter:  make(map[interface{}]func(interface{}, ...interface{})),
			CallRouter:  make(map[interface{}]func(interface{}, ...interface{}) *CallRet),
			dispatcher:  &jun_timer.Dispatcher{ChanTimer: make(chan *jun_timer.Timer, 10)},
			ModuleName:  name,
			Mod:         serv,
		}
		serv.SetModule(m)
		serv.SetServerName(name)
		mods.Store(name, m)

		m.Start(closeSig)
		fmt.Println("process exit:", name)
		mods.Delete(name)
		wg.Done()
	}()
	return nil
}

func StopAll() {
	StopAllServer()
	wg.Wait()
}

func StopAllServer() {
	mods.Range(func(key, value any) bool {
		ModeName := key.(string)
		fmt.Println("modName", ModeName)
		Stop(ModeName)
		return true
	})
}

type Module struct {
	ModuleName  string
	CallRouter  map[interface{}]func(interface{}, ...interface{}) *CallRet
	CastRouter  map[interface{}]func(interface{}, ...interface{})
	ChanCall    chan CallInfo
	ChanCallRet chan CallRet
	ChanCast    chan CastInfo
	ChanExit    chan ExitSig
	State       interface{}
	Mod         ModuleBehavior
	dispatcher  *jun_timer.Dispatcher
}

func (m *Module) RegisterCast(key interface{}, f func(interface{}, ...interface{})) {
	m.CastRouter[key] = f
}

func (m *Module) RegisterCall(key interface{}, f func(interface{}, ...interface{}) *CallRet) {
	m.CallRouter[key] = f
}

func (m *Module) HandlerCast(key interface{}, state interface{}, msg interface{}) {
	if _, ok := m.CastRouter[key]; ok {
		m.CastRouter[key](state, msg)
	} else {
		jun_log.Error("cast %s not found", key)
	}
}
func (m *Module) HandlerCall(key interface{}, state interface{}, msg interface{}) *CallRet {
	if _, ok := m.CallRouter[key]; ok {
		return m.CallRouter[key](state, msg)
	} else {
		jun_log.Error("call %s not found", key)
		return &CallRet{Error: fmt.Errorf("call %s not found", key)}
	}
}
func (m *Module) Start(closeSig chan ExitSig) {
	m.Mod.RegisterEvent()
	go m.Mod.Start(m.State)
	m.loop(closeSig)
}
func (m *Module) loop(closeSig chan ExitSig) {

	defer func() {
		if msg := recover(); msg != nil {

			fmt.Println("调用栈:%v", string(debug.Stack()), msg)
			fmt.Println("进程Panic退出，执行Terminate,退出原因")

			if closeSig != nil {
				closeSig <- ExitSig{Reason: ExitReasonPanic, Data: msg}
			}

			m.Mod.Terminate(m.State)
		} else {
			fmt.Println("进程退出：", m.State)
		}
	}()

	for {
		select {
		case callInfo := <-m.ChanCall:
			go func() {
				jun_log.Debug("callInfo:", callInfo)
				retInfo := m.HandlerCall(callInfo.Key, m.State, callInfo.msg)
				callInfo.replyChan <- *retInfo
			}()
		case castInfo := <-m.ChanCast:
			//m.Mod.HandlerCast(castInfo, m.State)
			m.HandlerCast(castInfo.Key, m.State, castInfo.Msg)
		case exitInfo := <-m.ChanExit:
			fmt.Println("进程退出，执行Terminate,退出原因:", m.State, exitInfo.Reason)
			m.Mod.Terminate(m.State)
			if closeSig != nil {
				closeSig <- exitInfo
			}
			fmt.Println("发送退出消息", m.State, exitInfo)

			return
		case t := <-m.dispatcher.ChanTimer:
			t.Cb()
		}

	}

}
func (s *Module) AfterFunc(d time.Duration, cb func()) *jun_timer.Timer {
	return s.dispatcher.AfterFunc(d, cb)
}
