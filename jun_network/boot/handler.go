package boot

import (
	"github.com/lijianjunljj/jungo/jun_network"
	"github.com/lijianjunljj/jungo/jun_server"
	"reflect"
)

type Handler struct {
	routerKey string
	Processor network.Processor
	Serv      jun_server.Server
	routerMap map[interface{}]interface{}
}

func NewHandler(routerKey string, processor network.Processor) Handler {
	hd := Handler{
		routerKey: routerKey,
		Processor: processor,
	}
	hd.routerMap = make(map[interface{}]interface{})
	return hd
}

func (that *Handler) Register(m interface{}, h interface{}) {
	that.routerMap[m] = h
}

func (that *Handler) RouterAll(serv jun_server.Server) {
	that.Serv = serv
	for k, v := range that.routerMap {
		that.router(k, v)
	}
}

func (that *Handler) router(m interface{}, h interface{}) {
	that.Processor.Register(m)
	that.Processor.SetRouter(m, that.routerKey)
	that.Serv.RegisterCast(reflect.TypeOf(m), h.(func(interface{}, ...interface{})))
}
