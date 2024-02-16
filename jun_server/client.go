package jun_server

import (
	"fmt"
	"time"
)

const (
	ExitReasonNormal = "normal"
	ExitReasonPanic  = "panic"
)

type ExitSig struct {
	Reason string
	Data   interface{}
}

type CallRet struct {
	Error  error
	Replay interface{}
}
type CallInfo struct {
	Id        string
	Key       string
	replyChan chan CallRet
	msg       interface{}
}
type CastInfo struct {
	Key interface{}
	Msg interface{}
}

func Call(distName, key string, msg interface{}) CallRet {
	distMod, ok := mods.Load(distName)
	if !ok {
		return CallRet{Error: fmt.Errorf("dist mod not started:%s", distMod)}
	} else {
		retInfo := callLocal(key, distMod.(*Module), msg)
		return retInfo
	}
}

func callLocal(key string, distMod *Module, msg interface{}) CallRet {
	callInfo := CallInfo{Key: key, replyChan: make(chan CallRet), msg: msg}
	distMod.ChanCall <- callInfo
	retInfo := <-callInfo.replyChan
	return retInfo
}

func Cast(distName, key interface{}, msg interface{}) {
	distMod, ok := mods.Load(distName)
	if !ok {
		return
	}
	distMod.(*Module).ChanCast <- CastInfo{Key: key, Msg: msg}
}

func Stop(distName string) {
	go func() {
		distMod, ok := mods.Load(distName)
		if !ok {
			return
		}
		distMod.(*Module).ChanExit <- ExitSig{Reason: ExitReasonNormal}
	}()
}

func Exit(distName, reason string, data interface{}) {
	go func() {
		distMod, ok := mods.Load(distName)
		if !ok {
			return
		}
		distMod.(*Module).ChanExit <- ExitSig{Reason: reason, Data: data}
	}()
}

func SendAfter(d time.Duration, distName, key string, msg interface{}) {
	distMod, ok := mods.Load(distName)
	if !ok {
		return
	}
	distMod.(*Module).AfterFunc(d, func() {
		Cast(distName, key, msg)
	})
}
