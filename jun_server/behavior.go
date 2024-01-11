package jun_server

type ModuleBehavior interface {
	Start(...interface{})
	HandlerCall(CallInfo, interface{})
	HandlerCast(CastInfo, interface{})
	Terminate(interface{})
}
