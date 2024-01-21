package jun_server

import (
	"fmt"
	"github.com/lijianjunljj/jungo/jun_log"
	"github.com/lijianjunljj/jungo/jun_timer"
	"runtime/debug"
	"sync"
	"time"
)

var mods sync.Map

var wg sync.WaitGroup

func Start(name string, mod ModuleBehavior, closeSig chan ExitSig, state interface{}) {
	wg.Add(1)
	//closeSig := make(chan ExitSig)
	go func() {

		_, ok := mods.Load(name)
		if ok {
			jun_log.Error("module has started")
			return
		}
		m := &Module{State: state,
			ChanCall:    make(chan CallInfo),
			ChanCast:    make(chan CastInfo, 100),
			ChanCallRet: make(chan CallRet),
			ChanExit:    make(chan ExitSig),
			CastRouter:  make(map[interface{}]func([]interface{})),
			dispatcher:  &jun_timer.Dispatcher{ChanTimer: make(chan *jun_timer.Timer, 10)},
			ModuleName:  name,
			Mod:         mod,
		}
		mods.Store(name, m)

		m.Start(closeSig)
		fmt.Println("process exit:", name)
		mods.Delete(name)
		wg.Done()
	}()
	//<-closeSig

}

func StopAll() {
	StopAllServer()
	wg.Wait()
}

func StopAllServer()  {
	mods.Range(func(key, value any) bool {
		ModeName := key.(string)
		fmt.Println("modName", ModeName)
		Stop(ModeName)
		return true
	})
}

type Module struct {
	ModuleName  string
	CallRouter  map[string]func(interface{}, interface{}) *CallRet
	CastRouter  map[interface{}]func([]interface{})
	ChanCall    chan CallInfo
	ChanCallRet chan CallRet
	ChanCast    chan CastInfo
	ChanExit    chan ExitSig
	State       interface{}
	Mod         ModuleBehavior
	dispatcher  *jun_timer.Dispatcher
}

func (m *Module) RegisterCast(key interface{}, f func(args []interface{})) {
	m.CastRouter[key] = f
}
func (m *Module) HandlerCast(key interface{}, state interface{}, msg interface{}) {
	if _, ok := m.CastRouter[key]; ok {
		m.CastRouter[key]([]interface{}{state, msg})
	}
}
func (m *Module) Start(closeSig chan ExitSig) {
	m.Mod.RegisterEvent(m)
	m.Mod.Start(m.State)
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
		//case callInfo := <-m.ChanCall:
		//m.Mod.HandlerCall(callInfo, m.State)
		case castInfo := <-m.ChanCast:
			//m.Mod.HandlerCast(castInfo, m.State)
			m.HandlerCast(castInfo.Key, m.State, castInfo.Msg)
		case exitInfo := <-m.ChanExit:
			fmt.Println("进程退出，执行Terminate,退出原因:", m.State, exitInfo.Reason)
			m.Mod.Terminate(m.State)
			if closeSig != nil {
				closeSig <- exitInfo
			}
			fmt.Println("发送退出消息",  m.State,exitInfo)

			return
		case t := <-m.dispatcher.ChanTimer:
			t.Cb()
		}

	}

}
func (s *Module) AfterFunc(d time.Duration, cb func()) *jun_timer.Timer {
	return s.dispatcher.AfterFunc(d, cb)
}
