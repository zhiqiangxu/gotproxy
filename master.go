package gotproxy

import (
	"log"
	"sync"
	"time"

	"github.com/zhiqiangxu/qrpc"
)

// Master starts Redirector and Proxy
type Master struct {
	redirector Redirector
	proxy      *Proxy
	wg         sync.WaitGroup
}

// NewMaster creates a Master
func NewMaster() *Master {
	return &Master{}
}

// Start the master
func (m *Master) Start(listeningPort uint16) {

	redirector := NewRedirector()
	proxy := NewProxy(redirector)

	qrpc.GoFunc(&m.wg, func() {
		proxy.ListenAndServe(listeningPort)
	})
	m.proxy = proxy

	time.Sleep(time.Second)

	if err := redirector.Start(listeningPort); err != nil {
		log.Panicf("failed to start redirector: %s", err)
	}

	log.Println("redirector started")
	m.redirector = redirector
}

// Stop master
func (m *Master) Stop() {
	log.Println("master stopping")
	if err := m.redirector.Stop(); err != nil {
		log.Panicf("redirector stop failed: %s", err)
	}
	log.Println("redirector exit")
	// TODO why hang when nil pointer deref ?
	m.proxy.Shutdown()

	m.wg.Wait()
}
