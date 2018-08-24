package gotproxy

import (
	"log"
	"net"

	"github.com/zhiqiangxu/gotproxy/darwin"
)

type defaultRedirector struct {
	cs *darwin.ControlSocket
}

// NewRedirector returns Redirector
func NewRedirector() Redirector {

	cs := darwin.NewControlSocket()
	if cs == nil {
		return nil
	}

	return &defaultRedirector{cs: cs}

}

func (dr *defaultRedirector) Start(listeningPort uint16) error {

	ok := dr.cs.StartRedirect(listeningPort)

	log.Println("StartRedirect", ok)
	if ok {
		return nil
	}

	return ErrRedirectFail
}

func (dr *defaultRedirector) Stop() error {

	ok := dr.cs.Close()
	log.Println("StopRedirect", ok)
	if ok {
		return nil
	}

	return ErrRedirectStopFail
}

func (dr *defaultRedirector) GetOriginalDst(net.Conn) string {
	return ""
}
