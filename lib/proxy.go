package lib

import (
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
)

const (
	bufferSize           = 16 * 1024
	errClosedNetworkConn = "use of closed network connection"
)

type ProxyConfig struct {
	To                net.Conn
	From              net.Conn
	ConnectionTimeout time.Duration
}

type Proxy struct {
	to                net.Conn
	from              net.Conn
	connectionTimeout time.Duration
	toStats           *IoStats
	fromStats         *IoStats
	statsInterrupt    chan struct{}
}

func NewProxy(cfg ProxyConfig) (proxy Proxy, err error) {
	if cfg.To == nil {
		err = errors.Errorf("'To' must not be nil")
		return
	}

	if cfg.From == nil {
		err = errors.Errorf("'From' must not be nil")
		return
	}

	proxy.toStats = &IoStats{}
	proxy.fromStats = &IoStats{}
	proxy.from = cfg.From
	proxy.to = cfg.To

	if cfg.ConnectionTimeout == 0 {
		proxy.connectionTimeout = 10 * time.Second
	}

	return
}

func (p *Proxy) Transfer() (err error) {
	var errChan = make(chan error, 1)

	go func() {
		err2 := copy(p.to, p.from, p.toStats)
		p.to.Close()
		p.from.Close()
		errChan <- err2
	}()

	err1 := copy(p.from, p.to, p.fromStats)
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

func copy(to io.Writer, from io.Reader, stats *IoStats) (err error) {
	var (
		buf    = make([]byte, bufferSize)
		readN  int
		writeN int
	)

	for {
		readN, err = from.Read(buf)
		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			break
		}

		if readN > 0 {
			stats.Rx += uint64(readN)
			writeN, err = to.Write(buf[0:readN])
			if err != nil {
				break
			}

			if readN != writeN {
				err = io.ErrShortWrite
				break
			}

			if writeN > 0 {
				stats.Tx += uint64(writeN)
			}
		}
	}

	if err != nil {
		e, _ := err.(*net.OpError)
		if e.Err.Error() == errClosedNetworkConn {
			err = nil
			return
		}
	}

	return
}
