package node

//
//type Caller struct {
//	jun_server.Server
//
//	client   *client.Client
//	nodeName string
//	distName string
//	callerId string
//}
//
//func Call(nodeName, distName string, client *client.Client) {
//	callerId := misc.GetUnixIDStr()
//	caller := &Caller{
//		callerId: callerId,
//		nodeName: nodeName,
//		distName: distName,
//		Server:   jun_server.Server{CloseSig: make(chan jun_server.ExitSig)},
//		client:   client,
//	}
//	jun_server.Run(callerId, caller, caller.CloseSig, callerId)
//
//	//retInfo := <-caller.CloseSig
//}
//
//func (that *Caller) Start(interface{}) {
//
//}
//func (that *Caller) RegisterEvent() {
//
//}
//func (that *Caller) Terminate(interface{}) {
//
//}
