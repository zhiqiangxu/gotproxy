package gotproxy

import (
	"errors"
	"net"
)

var (
	// ErrRedirectFail when redirect fail
	ErrRedirectFail = errors.New("failed to redirect")
	// ErrRedirectStopFail when fail to stop
	ErrRedirectStopFail = errors.New("failed to stop redirect")
)

// Redirector is responsible for redirect traffic to Proxy
type Redirector interface {
	Start(listeningPort uint16) error
	Stop() error

	GetOriginalDst(net.Conn) string
}
