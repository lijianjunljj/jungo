package jun_server

type ModuleBehavior interface {
	Start(interface{})
	RegisterEvent()
	GetServerName() string
	SetServerName(string)
	SetModule(*Module)
	GetState() interface{}
	GetCloseSig() chan ExitSig
	Terminate(interface{})
}
