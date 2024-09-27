package jun_node

import (
	"github.com/lijianjunljj/jungo/jun_node/client"
)

func Cast(nodeName, distName, key string, data interface{}) {
	caster := client.NewCaster(nodeName, distName, key, data)
	caster.Cast()
	return
}
