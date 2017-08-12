package lib

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// DumpTcpServer is an echo server which copies
// what is received in its input to a 'writer' and
// then writes back what comes to that connection
// that was initiated.
type DumbTcpServer struct {
	w    io.Writer
	ln   net.Listener
	port int
}

func NewDumbTcpServer(w io.Writer) DumbTcpServer {
	if w == nil {
		panic(fmt.Errorf("a non-nil writer must be passed"))
	}

	return DumbTcpServer{w: w}
}

func (s *DumbTcpServer) Listen() (err error) {
	ln, err := net.Listen("tcp4", ":0")
	defer ln.Close()
	if err != nil {
		return
	}

	s.ln = ln
	s.port = ln.Addr().(*net.TCPAddr).Port

	for {
		conn, err := ln.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				continue
			}
			fmt.Printf("listen errored: %+v\n", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *DumbTcpServer) GetPort() int {
	return s.port
}

func (s *DumbTcpServer) Close() {
	if s.ln != nil {
		s.ln.Close()
	}
}

func (s *DumbTcpServer) handle(conn net.Conn) {
	var err error

	defer conn.Close()
	_, err = io.Copy(conn, io.TeeReader(conn, s.w))
	if err != nil {
		fmt.Println("ERRORED while echoing", err)
	}
}

func TestDumbTcpServerTransmitsData(t *testing.T) {
	var msg = []byte("PING\r\n")
	var receiveBuffer = make([]byte, len(msg))
	var buf bytes.Buffer

	server := NewDumbTcpServer(&buf)
	defer server.Close()

	go func() {
		assert.NoError(t, server.Listen())
	}()

	time.Sleep(100 * time.Millisecond)

	addr := fmt.Sprintf("localhost:%d", server.GetPort())

	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	defer conn.Close()

	n, err := conn.Write(msg)
	assert.NoError(t, err)
	assert.Equal(t, len(msg), n)

	n, err = conn.Read(receiveBuffer)
	assert.NoError(t, err)
	assert.Equal(t, len(msg), n)
	assert.Equal(t, msg, buf.Bytes())
}
