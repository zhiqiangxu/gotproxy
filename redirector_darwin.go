package gotproxy

import (
	"encoding/binary"
	"io"
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

func (dr *defaultRedirector) GetOriginalDst(rw net.Conn) (string, uint16, error) {

	bytes := make([]byte, 1)
	_, err := io.ReadFull(rw, bytes)
	if err != nil {
		return "", 0, err
	}
	length := uint8(bytes[0])
	bytes = make([]byte, length)
	_, err = io.ReadFull(rw, bytes)
	if err != nil {
		return "", 0, err
	}
	dst := string(bytes)
	bytes = make([]byte, 2)
	_, err = io.ReadFull(rw, bytes)
	if err != nil {
		return "", 0, err
	}
	port := binary.BigEndian.Uint16(bytes)

	return dst, port, nil
}
