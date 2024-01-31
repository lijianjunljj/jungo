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
	Id        string
	Key       string
	replyChan chan CallRet
	msg       interface{}
}
type CastInfo struct {
	Key interface{}
	Msg interface{}
}

func SyncCall(distName string, msg interface{}) CallRet {
	return CallLocal(distName, msg)
}

func AsyncCall(distName string, msg interface{}, f func(CallRet)) {
	go func() {
		retInfo := CallLocal(distName, msg)
		f(retInfo)
	}()
}

func CallLocal(distName string, msg interface{}) CallRet {
	distMod, ok := mods.Load(distName)
	if !ok {
		return CallRet{err: fmt.Errorf("dist mod not started:%s", distMod)}
	}
	callInfo := CallInfo{replyChan: make(chan CallRet), msg: msg}
	distMod.(*Module).ChanCall <- callInfo
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
		fmt.Println("distMod, ok:", distMod, ok)
		if !ok {
			return
		}
		distMod.(*Module).ChanExit <- ExitSig{Reason: ExitReasonNormal}
		fmt.Println("sssssssssssssssssss")
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
