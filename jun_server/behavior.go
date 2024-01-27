package jun_server

type ModuleBehavior interface {
	Start(...interface{})
	RegisterEvent()
	GetServerName() string
	SetModule(*Module)
	GetState() interface{}
	GetCloseSig() chan ExitSig
	Terminate(interface{})
}
