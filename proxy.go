package gotproxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/zhiqiangxu/qrpc"
)

// Proxy deals with redirected traffic
type Proxy struct {
	inconns    sync.Map
	redirector Redirector
	ln         net.Listener
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewProxy returns a Proxy instance
func NewProxy(redirector Redirector) *Proxy {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Proxy{redirector: redirector, ctx: ctx, cancelFunc: cancelFunc}
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

	p.inconns.Store(rw, struct{}{})

	defer func() {
		rw.Close()
		p.inconns.Delete(rw)
		if err := recover(); err != nil {
			log.Println("panic", err)
		}
	}()

	dst, port, err := p.redirector.GetOriginalDst(rw)
	if err != nil {
		log.Println("GetOriginalDst err", err)
		return
	}

	targetAddr := fmt.Sprintf("%s:%d", dst, port)
	p.forward(rw, targetAddr)
}

func (p *Proxy) forward(rw net.Conn, targetAddr string) {
	log.Println("new conn")
	conn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Println("Dial err", err)
		return
	}

	var wg sync.WaitGroup

	qrpc.GoFunc(&wg, func() {
		io.Copy(conn, rw)
	})
	qrpc.GoFunc(&wg, func() {
		io.Copy(rw, conn)
	})

	wg.Wait()
}

// Shutdown stops the Proxy
func (p *Proxy) Shutdown() error {
	p.cancelFunc()
	p.ln.Close()

	p.inconns.Range(func(k, v interface{}) bool {
		k.(net.Conn).Close()
		return true
	})
	p.wg.Wait()
	return nil
}
