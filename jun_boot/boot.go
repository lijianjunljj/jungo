package jun_boot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lijianjunljj/jungo/jun_log"
	"github.com/lijianjunljj/jungo/jun_server"
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
	// 监听系统退出信号，在 Linux 上支持 kill 发送的 SIGTERM
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-c
	fmt.Println("接收到退出信号:", sig)
	jun_server.StopAll()
	fmt.Println("协程全部退出")
	signal.Stop(c)
}
