package server

import "github.com/lijianjunljj/jungo/jun_network/json"

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&LoginTos{})
	Processor.Register(&LoginToc{})
}

type MsgComToc struct {
	Code int8   `json:"code"`
	Msg  string `json:"msg"`
}
type LoginTos struct {
	Cookie   string
	NodeName string
}

type LoginToc struct {
	MsgComToc
}
