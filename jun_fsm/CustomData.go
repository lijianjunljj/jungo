package jun_fsm

type CustomData struct {
	CallFuncOption ICustomDataFuncOptions
	Data           interface{}
	Fsm            *Fsm
	State          *State
}
