package jun_node

import (
	"github.com/lijianjunljj/jungo/jun_node/client"
	"github.com/lijianjunljj/jungo/jun_server"
)

func Cast(nodeName, distName, key string, data interface{}) {
	caller := client.NewCaster(nodeName, distName, key, data)
	err := jun_server.RunServer(caller)
	if err != nil {
		return
	}
	return
}
