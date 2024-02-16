package jun_node

import (
	"errors"
	"github.com/lijianjunljj/jungo/jun_node/client"
	"github.com/lijianjunljj/jungo/jun_server"
)

func Call(nodeName, distName, key string, data interface{}) (ret *jun_server.CallRet) {
	ret = new(jun_server.CallRet)
	caller := client.NewCaller(nodeName, distName, key, data)
	err := jun_server.RunServer(caller)
	if err != nil {
		ret.Error = err
		return
	}
	closeInfo := <-caller.CloseSig
	ret = &jun_server.CallRet{Error: errors.New(closeInfo.Reason), Replay: closeInfo.Data}
	return
}
