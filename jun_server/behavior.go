package jun_server

type ModuleBehavior interface {
	Start(...interface{})
	RegisterEvent(*Module)
	//HandlerCall(CallInfo, interface{})
	//HandlerCast(CastInfo, interface{})
	Terminate(interface{})
}
