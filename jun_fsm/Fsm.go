// Package jun_fsm author: lijianjun qq:853944509
package jun_fsm

import (
	"fmt"
	"github.com/lijianjunljj/jungo/jun_server"
	"sync"
	"time"
)

type ExitSig struct {
	Reason string
	Data   interface{}
}

func NewFsm(data *CustomData, states []*State) *Fsm {
	data.State = states[0]
	return &Fsm{
		loopState: &LoopState{customData: data},
		states:    states,
	}
}

type Fsm struct {
	jun_server.ModuleBehavior
	lock          sync.RWMutex
	loopState     *LoopState
	currSateIndex int
	states        []*State
	oldStates     []*State
	isInit        bool
}

func (that *Fsm) Exit() {
	fmt.Println("房间进程退出.....")
	//that.loopState.Cancel()
}
func (that *Fsm) GetCurrentState() *State {
	newState := that.states[that.currSateIndex]
	return newState
}

func (that *Fsm) RemoveState(name string) {
	var newStates []*State
	currIndex := 0
	for _, v := range that.states {
		if v.Name != name {
			if that.states[that.currSateIndex].Name == v.Name {
				that.currSateIndex = currIndex
			}
			newStates = append(newStates, v)
			currIndex++
		}
	}
	that.states = newStates
}

func (that *Fsm) NextState() {
	state := that.states[that.currSateIndex]
	if that.currSateIndex >= len(that.states)-1 {
		that.currSateIndex = 0
	} else {
		that.currSateIndex++
	}
	if state.hookOption.EndFunc != nil {
		state.hookOption.EndFunc(that, state, that.loopState.customData.Data)
	}

	newState := that.states[that.currSateIndex]
	that.loopState.customData.State = newState
	newState.Reset()
	if newState.hookOption.StartFunc != nil {
		newState.hookOption.StartFunc(that, newState, that.loopState.customData.Data)
	}
	if that.loopState.customData != nil {
		that.loopState.customData.CallFuncOption.StatusChangeFunc(that, state, newState, that.loopState.customData.Data)
	}
}
func (that *Fsm) SwitchStates(states []*State) {
	that.oldStates = []*State{}
	if that.states == nil {
		for _, v := range that.states {
			that.oldStates = append(that.oldStates, v)
		}
	}
	that.currSateIndex = 0
	that.states = states
}

func (that *Fsm) Start(arg ...interface{}) {
	fmt.Println("start............")
	roomId := arg[0].(string)
	jun_server.SendAfter(1*time.Second, roomId, "loop")
}

func (s *Fsm) HandlerCall(callInfo jun_server.CallInfo, data interface{}) {
	fmt.Println("HandlerCall:", callInfo, data)

}
func (s *Fsm) HandlerCast(castInfo jun_server.CastInfo, data interface{}) {
	if castInfo.Msg == "loop" {
		jun_server.SendAfter(1*time.Second, data.(string), "loop")
		s.Loop()
	}

}
func (s *Fsm) Terminate(data interface{}) {
	fmt.Println("Terminate:", data)
}
func (that *Fsm) Loop() {
	if len(that.states) > that.currSateIndex {
		state := that.states[that.currSateIndex]
		//fmt.Println("leftTime:", state)
		if !that.isInit {
			that.isInit = true
			state.hookOption.StartFunc(that, state, that.loopState.customData.Data)
		}

		if state.LeftTime == -1 {
			if state.hookOption.LoopFunc != nil {
				state.hookOption.LoopFunc(that, state, that.loopState.customData.Data)
			}
		} else if state.LeftTime <= 0 {
			that.NextState()
		} else {
			if state.hookOption.LoopFunc != nil {
				state.hookOption.LoopFunc(that, state, that.loopState.customData.Data)
			}
			state.LeftTime--
			if that.loopState.customData != nil {
				that.loopState.customData.CallFuncOption.LoopFunc(that, state, that.loopState.customData.Data)
				that.loopState.customData.State = that.states[that.currSateIndex]
			}
		}

	}
}
