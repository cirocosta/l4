package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"
)

const (
	bufferSize = 16 * 1024
)

func NewProxy(cfg ProxyConfig) (proxy Proxy, err error) {
	proxy.toStats = &IoStats{}
	proxy.fromStats = &IoStats{}
	proxy.from = cfg.From
	proxy.to = cfg.To
	proxy.connectionTimeout = cfg.ConnectionTimeout

	return
}

func (c *GracefulConn) Close() error {
	err := c.Conn.Close()
	if err != nil {
		return err
	}

	c.ln.closeConn()
	return nil
}

func NewGracefulListener(cfg GracefulListenerConfig) net.Listener {
	return &GracefulListener{
		ln:               cfg.Listener,
		maxCloseWaitTime: cfg.MaximumWaitTime,
		done:             make(chan struct{}),
	}
}

func (ln *GracefulListener) Accept() (net.Conn, error) {
	c, err := ln.ln.Accept()
	if err != nil {
		return nil, err
	}

	atomic.AddUint64(&ln.connsCount, 1)
	return &GracefulConn{
		Conn: c,
		ln:   ln,
	}, nil
}

func (ln *GracefulListener) Addr() net.Addr {
	return ln.ln.Addr()
}

func (ln *GracefulListener) Close() error {
	err := ln.ln.Close()
	if err != nil {
		return err
	}

	return ln.waitForZeroConns()
}

func (ln *GracefulListener) waitForZeroConns() error {
	atomic.AddUint64(&ln.shutdown, 1)
	if atomic.LoadUint64(&ln.connsCount) == 0 {
		close(ln.done)
		return nil
	}

	select {
	case <-ln.done:
		return nil
	case <-time.After(ln.maxCloseWaitTime):
		return fmt.Errorf("cannot complete graceful shutdown in %s",
			ln.maxCloseWaitTime)
	}

	return nil
}

func (ln *GracefulListener) closeConn() {
	connsCount := atomic.AddUint64(&ln.connsCount, ^uint64(0))
	if atomic.LoadUint64(&ln.shutdown) != 0 && connsCount == 0 {
		close(ln.done)
	}
}

func (p *Proxy) Transfer() (err error) {
	var errChan = make(chan error, 1)

	go func() {
		err2 := Copy(p.to, p.from, p.toStats)
		p.to.Close()
		p.from.Close()
		errChan <- err2
	}()

	err1 := Copy(p.from, p.to, p.fromStats)
	p.to.Close()
	p.from.Close()
	err2 := <-errChan

	if err1 != nil {
		err = err1
		return
	}

	err = err2
	return
}

func handleConnection(conn net.Conn) {
	agent, err := net.Dial("tcp4", "0.0.0.0:8080")
	if err != nil {
		log.Panicf("couldn't dial port 8080 - %+v\n", err)
	}

	proxy, err := NewProxy(ProxyConfig{
		To:                agent,
		From:              conn,
		ConnectionTimeout: 10 * time.Second,
	})
	if err != nil {
		log.Panicf("couldn't create proxy - %+v\n", err)
	}

	err = proxy.Transfer()
	if err != nil {
		log.Panicf("errored transfering between connections - %+v\n", err)
	}
}

func Copy(to io.Writer, from io.Reader, stats *IoStats) (err error) {
	var (
		buf    = make([]byte, bufferSize)
		readN  int
		writeN int
	)

	fmt.Println("%+v\n", to)
	fmt.Println("%+v\n", from)
	fmt.Println("%+v\n", stats)

	for {
		readN, err = from.Read(buf)
		if err != nil {
			return
		}

		if err == io.EOF {
			return
		}

		if readN > 0 {
			stats.Rx += uint64(readN)
			writeN, err = to.Write(buf[0:readN])
			if err != nil {
				return
			}

			if readN != writeN {
				err = io.ErrShortWrite
				return
			}

			if writeN > 0 {
				stats.Tx += uint64(writeN)
			}
		}
	}

	return
}

func main() {
	ln, err := net.Listen("tcp4", ":8000")
	if err != nil {
		log.Panicf("couldn't listen on port 8000", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panicf("errored accepting connection", err)
		}

		log.Printf("connection accepted [local=%s,remote=%s]\n",
			conn.LocalAddr().String(), conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}
