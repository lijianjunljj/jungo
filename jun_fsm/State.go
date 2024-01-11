package jun_fsm

type AbstractState interface {
	GetStartFunc() func(*Fsm, *State, interface{})
	GetLoopFunc() func(*Fsm, *State, interface{})
	GetEndFunc() func(*Fsm, *State, interface{})
	GetStatusChangeFunc() func(*Fsm, *State, int, int)
}

type StateFuncOption interface {
	StartFunc(*Fsm, *State, interface{})
	LoopFunc(*Fsm, *State, interface{})
	EndFunc(*Fsm, *State, interface{})
	StatusChangeFunc(*Fsm, *State, int, int)
}

type State struct {
	AbstractState
	hookOption StateFuncOption
	Name       string
	LeftTime   int
	RightTime  int
	TotalTime  int
}

func NewFsmState(name string, totalTime int, option StateFuncOption) *State {
	state := &State{
		Name:      name,
		LeftTime:  totalTime,
		TotalTime: totalTime,
	}
	state.hookOption = option
	return state
}
func (that *State) Reset() {
	that.LeftTime = that.TotalTime
}
