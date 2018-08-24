package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/zhiqiangxu/gotproxy"
)

const (
	listeningPort = 9999
)

func main() {

	master := gotproxy.NewMaster()
	master.Start(listeningPort)

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-quitChan

	master.Stop()
}
