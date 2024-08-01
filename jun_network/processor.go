package network

type Processor interface {
	Register(interface{}) string
	SetRouter(interface{}, string)
	// must goroutine safe
	Route(msg interface{}, userData interface{}) error
	// must goroutine safe
	Unmarshal(data []byte) (interface{}, error)
	// must goroutine safe
	Marshal(msg interface{}) ([][]byte, error)
}
