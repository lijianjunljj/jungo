package jun_server

import (
	"fmt"
	"time"
)

func TestGenServer() {

	//closeSing := make(chan ExitSig)
	//go Start("server", &server{}, closeSing, "123456")
	////exitInfo := <-closeSing
	////fmt.Println("exitInfo:", exitInfo)
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt, os.Kill)
	//sig := <-c
	//fmt.Println("sig:", sig)
}

type server struct {
	ModuleBehavior
}

func (s *server) Start(args ...interface{}) {
	data := args[0]
	fmt.Println("start:", data)
	SendAfter(1*time.Second, "server", "loop")
}
func (s *server) HandlerCall(callInfo CallInfo, data interface{}) {
	fmt.Println("HandlerCall:", callInfo, data)

}
func (s *server) HandlerCast(castInfo CastInfo, data interface{}) {
	fmt.Println("HandlerCast:", castInfo, data)
	SendAfter(1*time.Second, "server", "loop")
}
func (s *server) Terminate(data interface{}) {
	fmt.Println("Terminate:", data)
}
