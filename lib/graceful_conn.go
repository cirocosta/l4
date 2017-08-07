package lib

import (
	"net"
)

type GracefulConn struct {
	net.Conn
	ln *GracefulListener
}

func (c *GracefulConn) Close() error {
	err := c.Conn.Close()
	if err != nil {
		return err
	}

	c.ln.closeConn()
	return nil
}
