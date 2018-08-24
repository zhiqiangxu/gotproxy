package gotproxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/zhiqiangxu/qrpc"
)

// Proxy deals with redirected traffic
type Proxy struct {
	redirector Redirector
	ln         net.Listener
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewProxy returns a Proxy instance
func NewProxy(redirector Redirector) *Proxy {
	ctx, cancelCtx := context.WithCancel(context.Background())
	return &Proxy{redirector: redirector, ctx: ctx, cancelFunc: cancelCtx}
}

// ListenAndServe starts listening
func (p *Proxy) ListenAndServe(listeningPort uint16) error {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", listeningPort))
	if err != nil {
		return err
	}
	p.ln = ln

	p.wg.Add(1)
	defer p.wg.Done()

	for {
		rw, err := ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				time.Sleep(time.Second)
				continue
			}
			return err
		}

		qrpc.GoFunc(&p.wg, func() {
			p.serve(rw)
		})

		select {
		case <-p.ctx.Done():
			return p.ctx.Err()
		default:
		}
	}
}

func (p *Proxy) serve(rw net.Conn) {
	defer rw.Close()

	dst := p.redirector.GetOriginalDst(rw)
	p.forward(rw, dst)
}

func (p *Proxy) forward(rw net.Conn, dst string) {
	log.Println("new conn")
}

// Shutdown stops the Proxy
func (p *Proxy) Shutdown() error {
	p.cancelFunc()
	p.ln.Close()
	p.wg.Wait()
	return nil
}
