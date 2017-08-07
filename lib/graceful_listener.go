package lib

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultMaxCloseWaitTime = 5 * time.Second
)

type GracefulListenerConfig struct {
	Listener        net.Listener
	MaximumWaitTime time.Duration
}

type GracefulListener struct {
	ln               net.Listener
	maxCloseWaitTime time.Duration
	done             chan struct{}
	connsCount       uint64
	shutdown         uint64
}

func NewGracefulListener(cfg GracefulListenerConfig) net.Listener {
	if cfg.Listener == nil {
		panic(errors.New(
			"Can't create graceful listener without a listener"))
	}

	if cfg.MaximumWaitTime == 0 {
		cfg.MaximumWaitTime = defaultMaxCloseWaitTime
	}

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
		return errors.Errorf("cannot complete graceful shutdown in %s",
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
