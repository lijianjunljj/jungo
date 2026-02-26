// Package jun_fsm author: lijianjun qq:853944509
package jun_fsm

import (
	"fmt"
	"sync"
	"time"

	"github.com/lijianjunljj/jungo/jun_server"
)

type ExitSig struct {
	Reason string
	Data   interface{}
}

func NewFsm(data *CustomData, states []*State) *Fsm {
	data.State = states[0]
	return &Fsm{
		skipStateIndex: -1,
		loopState:      &LoopState{customData: data},
		states:         states,
	}
}

type Fsm struct {
	jun_server.Server
	lock           sync.RWMutex
	loopState      *LoopState
	currSateIndex  int
	states         []*State
	oldStates      []*State
	isInit         bool
	skipStateIndex int
}

func (that *Fsm) Exit() {
	fmt.Println("房间进程退出.....")
	//that.loopState.Cancel()
}
func (that *Fsm) GetCurrentState() *State {
	newState := that.states[that.currSateIndex]
	return newState
}
func (that *Fsm) UnshiftState(state *State) {
	that.states = append([]*State{state}, that.states...)
	that.currSateIndex++
}

func (that *Fsm) RemoveState(name string) {
	that.lock.Lock()
	defer that.lock.Unlock()

	// 检查要删除的状态是否存在，并记录其索引
	targetIndex := -1
	for i, v := range that.states {
		if v.Name == name {
			targetIndex = i
			break
		}
	}

	// 如果状态不存在，直接返回
	if targetIndex == -1 {
		return
	}

	// 构建新的状态列表
	var newStates []*State
	for _, v := range that.states {
		if v.Name != name {
			newStates = append(newStates, v)
		}
	}

	// 调整当前状态索引
	if targetIndex < that.currSateIndex {
		// 删除的是当前状态之前的状态，currSateIndex 需要减1
		that.currSateIndex--
	} else if targetIndex == that.currSateIndex {
		// 删除的是当前状态
		if len(newStates) == 0 {
			// 没有剩余状态了
			that.currSateIndex = 0
		} else {
			// 删除当前状态后，后面的状态会前移
			// 如果原来有后续状态，索引不变（指向原来的下一个状态）
			// 如果没有后续状态，索引需要减1（指向前一个状态）
			if that.currSateIndex >= len(newStates) {
				that.currSateIndex = len(newStates) - 1
			}
		}
	}
	// 如果 targetIndex > that.currSateIndex，删除的是当前状态之后的状态，currSateIndex 不变

	// 确保索引有效
	if len(newStates) > 0 {
		if that.currSateIndex < 0 {
			that.currSateIndex = 0
		} else if that.currSateIndex >= len(newStates) {
			that.currSateIndex = len(newStates) - 1
		}
	} else {
		that.currSateIndex = 0
	}

	that.states = newStates
}

func (that *Fsm) nextState() {
	state := that.states[that.currSateIndex]

	if state.hookOption != nil {
		state.hookOption.EndFunc(that, state, that.loopState.customData.Data)
	}

	var newState *State
	if that.skipStateIndex > -1 {
		that.currSateIndex = that.skipStateIndex
		that.skipStateIndex = -1
	} else {
		if that.currSateIndex >= len(that.states)-1 {
			that.currSateIndex = 0
		} else {
			that.currSateIndex++
		}
	}
	newState = that.states[that.currSateIndex]
	that.loopState.customData.State = newState
	newState.Reset()
	if newState.hookOption != nil {
		newState.hookOption.StartFunc(that, newState, that.loopState.customData.Data)
	}
	if that.loopState.customData != nil {
		that.loopState.customData.CallFuncOption.StatusChangeFunc(that, state, newState, that.loopState.customData.Data)
	}
}

func (that *Fsm) Start(arg interface{}) {
	fmt.Println("start............")
	roomId := arg.(string)
	jun_server.SendAfter(1*time.Second, roomId, "loop", nil)
}
func (s *Fsm) RegisterEvent() {
	s.RegisterCast("loop", s.HandlerLoop)
}
func (s *Fsm) HandlerLoop(state interface{}, args ...interface{}) {
	data := args[0]
	jun_server.SendAfter(1*time.Second, state.(string), "loop", data)
	s.loop()
}

func (s *Fsm) Terminate(data interface{}) {
	fmt.Println("Fsm Terminate:", data)
}

// JumpToNextState 立即跳转到下一个状态（带锁保护）
// skipEndFunc: 如果为 true，则跳过当前状态的 EndFunc 调用（用于在 StartFunc 中调用时避免递归）
func (that *Fsm) JumpToNextState(skipEndFunc ...bool) error {

	if len(that.states) == 0 {
		return fmt.Errorf("没有可用状态")
	}

	// 获取当前状态
	currentState := that.states[that.currSateIndex]

	// 判断是否跳过 EndFunc（默认不跳过）
	shouldSkipEndFunc := false
	if len(skipEndFunc) > 0 {
		shouldSkipEndFunc = skipEndFunc[0]
	}

	// 执行当前状态的EndFunc（除非明确要求跳过）
	if !shouldSkipEndFunc && currentState.hookOption != nil {
		currentState.hookOption.EndFunc(that, currentState, that.loopState.customData.Data)
	}

	// 计算下一个状态索引
	var nextIndex int
	if that.currSateIndex >= len(that.states)-1 {
		nextIndex = 0
	} else {
		nextIndex = that.currSateIndex + 1
	}

	// 切换到下一个状态
	that.currSateIndex = nextIndex
	nextState := that.states[that.currSateIndex]
	that.loopState.customData.State = nextState
	nextState.Reset()

	// 执行下一个状态的StartFunc
	if nextState.hookOption != nil {
		nextState.hookOption.StartFunc(that, nextState, that.loopState.customData.Data)
	}

	// 触发状态变更回调
	if that.loopState.customData != nil {
		that.loopState.customData.CallFuncOption.StatusChangeFunc(that, currentState, nextState, that.loopState.customData.Data)
	}

	return nil
}

// JumpToState 立即跳转到指定名称的状态（带锁保护）
// skipEndFunc: 如果为 true，则跳过当前状态的 EndFunc 调用（用于在 StartFunc 中调用时避免递归）
func (that *Fsm) JumpToState(name string, skipEndFunc ...bool) error {
	fmt.Println("JumpToState:", name)

	fmt.Println("11111111111")
	if len(that.states) == 0 {
		return fmt.Errorf("没有可用状态")
	}

	// 查找目标状态
	targetIndex := -1
	for i, v := range that.states {
		if v.Name == name {
			targetIndex = i
			break
		}
	}
	fmt.Println("targetIndex:", targetIndex)

	if targetIndex == -1 {
		return fmt.Errorf("状态 %s 不存在", name)
	}

	// 如果已经是目标状态，直接返回
	if that.currSateIndex == targetIndex {
		return nil
	}

	// 获取当前状态
	currentState := that.states[that.currSateIndex]
	fmt.Println("currentState:", currentState.Name)

	// 判断是否跳过 EndFunc（默认不跳过）
	shouldSkipEndFunc := false
	if len(skipEndFunc) > 0 {
		shouldSkipEndFunc = skipEndFunc[0]
	}

	// 执行当前状态的EndFunc（除非明确要求跳过）
	if !shouldSkipEndFunc && currentState.hookOption != nil {
		currentState.hookOption.EndFunc(that, currentState, that.loopState.customData.Data)
	}

	// 切换到目标状态
	that.currSateIndex = targetIndex
	targetState := that.states[that.currSateIndex]
	that.loopState.customData.State = targetState
	targetState.Reset()

	// 执行目标状态的StartFunc
	if targetState.hookOption != nil {
		targetState.hookOption.StartFunc(that, targetState, that.loopState.customData.Data)
	}

	// 触发状态变更回调
	if that.loopState.customData != nil {
		that.loopState.customData.CallFuncOption.StatusChangeFunc(that, currentState, targetState, that.loopState.customData.Data)
	}

	return nil
}

// JumpToStateByIndex 立即跳转到指定索引的状态（带锁保护）
// skipEndFunc: 如果为 true，则跳过当前状态的 EndFunc 调用（用于在 StartFunc 中调用时避免递归）
func (that *Fsm) JumpToStateByIndex(index int, skipEndFunc ...bool) error {

	if len(that.states) == 0 {
		return fmt.Errorf("没有可用状态")
	}

	if index < 0 || index >= len(that.states) {
		return fmt.Errorf("状态索引 %d 超出范围 [0, %d)", index, len(that.states))
	}

	// 如果已经是目标状态，直接返回
	if that.currSateIndex == index {
		return nil
	}

	// 获取当前状态
	currentState := that.states[that.currSateIndex]

	// 判断是否跳过 EndFunc（默认不跳过）
	shouldSkipEndFunc := false
	if len(skipEndFunc) > 0 {
		shouldSkipEndFunc = skipEndFunc[0]
	}

	// 执行当前状态的EndFunc（除非明确要求跳过）
	if !shouldSkipEndFunc && currentState.hookOption != nil {
		currentState.hookOption.EndFunc(that, currentState, that.loopState.customData.Data)
	}

	// 切换到目标状态
	that.currSateIndex = index
	targetState := that.states[that.currSateIndex]
	that.loopState.customData.State = targetState
	targetState.Reset()

	// 执行目标状态的StartFunc
	if targetState.hookOption != nil {
		targetState.hookOption.StartFunc(that, targetState, that.loopState.customData.Data)
	}

	// 触发状态变更回调
	if that.loopState.customData != nil {
		that.loopState.customData.CallFuncOption.StatusChangeFunc(that, currentState, targetState, that.loopState.customData.Data)
	}

	return nil
}
func (that *Fsm) loop() {
	if len(that.states) > that.currSateIndex {
		state := that.states[that.currSateIndex]
		//fmt.Println("leftTime:", state)
		if !that.isInit {
			that.isInit = true
			if state.hookOption != nil {
				state.hookOption.StartFunc(that, state, that.loopState.customData.Data)
			}
		}

		if state.LeftTime == -1 {
			if state.hookOption != nil {
				state.hookOption.LoopFunc(that, state, that.loopState.customData.Data)
			}
		} else if state.LeftTime <= 0 {
			that.nextState()
		} else {
			if state.hookOption != nil {
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
