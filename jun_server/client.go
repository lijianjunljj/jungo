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
	err error
	ret interface{}
}
type CallInfo struct {
	srcName string
	msg     interface{}
}
type CastInfo struct {
	Msg interface{}
}

func Call(srcName, distName string, msg interface{}) CallRet {
	distMod, ok := mods.Load(distName)
	if !ok {
		return CallRet{err: fmt.Errorf("dist mod not started:%s", distMod)}
	}
	srcMod, ok := mods.Load(srcName)
	if !ok {
		return CallRet{err: fmt.Errorf("src dist mod not started:%s", srcMod)}
	}
	distMod.(*Module).ChanCall <- CallInfo{srcName: srcName, msg: msg}
	retInfo := <-srcMod.(*Module).ChanCallRet
	return retInfo
}

func Cast(distName string, msg interface{}) {
	distMod, ok := mods.Load(distName)
	if !ok {
		return
	}
	distMod.(*Module).ChanCast <- CastInfo{Msg: msg}
}

func Stop(distName string) {
	distMod, ok := mods.Load(distName)
	fmt.Println("distMod, ok:", distMod, ok)
	if !ok {
		return
	}
	distMod.(*Module).ChanExit <- ExitSig{Reason: ExitReasonNormal}
	fmt.Println("sssssssssssssssssss")
	//mods.Delete(distName)
}

func SendAfter(d time.Duration, distName string, msg interface{}) {
	distMod, ok := mods.Load(distName)
	if !ok {
		return
	}
	distMod.(*Module).AfterFunc(d, func() {
		Cast(distName, msg)
	})
}
