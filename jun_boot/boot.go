package jun_boot

import (
	"fmt"
	"github.com/lijianjunljj/jungo/jun_log"
	"github.com/lijianjunljj/jungo/jun_server"
	"os"
	"os/signal"
)

func Start() {
	jun_log.Init()
}
func Run(f func()) {
	jun_log.Init()
	f()
	End()
}
func End() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("接收到退出信号:", sig)
	jun_server.StopAll()
	fmt.Println("协程全部退出")
}
