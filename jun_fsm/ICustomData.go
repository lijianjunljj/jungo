package jun_fsm

type ICustomDataFuncOptions interface {
	StatusChangeFunc(fsm *Fsm, oldState *State, newState *State, customData interface{})
	LoopFunc(fsm *Fsm, state *State, customData interface{})
	Init(string, interface{}) error
}
