package lib

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var testCases = []struct {
		description string
		config      *ProxyConfig
		shouldError bool
	}{
		{
			description: "error on empty cfg",
			config:      &ProxyConfig{},
			shouldError: true,
		},
		{
			description: "succeed if all set",
			config: &ProxyConfig{
				To:                &net.TCPConn{},
				From:              &net.TCPConn{},
				ConnectionTimeout: 1 * time.Second,
			},
			shouldError: false,
		},
		{
			description: "fail if 'to' not set",
			config: &ProxyConfig{
				From:              &net.TCPConn{},
				ConnectionTimeout: 1 * time.Second,
			},
			shouldError: true,
		},
		{
			description: "fail if 'from' not set",
			config: &ProxyConfig{
				To:                &net.TCPConn{},
				ConnectionTimeout: 1 * time.Second,
			},
			shouldError: true,
		},
		{
			description: "doesnt fail if connection timeout not set",
			config: &ProxyConfig{
				To:   &net.TCPConn{},
				From: &net.TCPConn{},
			},
			shouldError: false,
		},
	}

	var (
		proxy Proxy
		err   error
	)

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			proxy, err = NewProxy(*tc.config)
			if tc.shouldError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

// Typically 'proxy' sits in the middle
// of two connections:
//	    s1    s2   s3    s4
//	conn1 <--> proxy <--> conn2
//
// such that:
//	- a write from 'conn1' to 'proxy'
//	becomes a write from 'proxy' to 'conn2'.
//
//	conn1:	s1.write()
//	proxy:	s2.read()
//	proxy: s3.write()
//	conn2:	s4.read()
//
//
//	- a write from 'conn2' to 'proxy'
//	becomes a write from 'proxy' to 'conn1'.
//
//	(example above but in oposite direction)
//
//
// Here we do something atypical that makes
// the proxy get the connections in a loop
// scenario due to the fact that there's an
// echo server in the middle.
//
//         s1       s2
//	conn1 <--> echo
//
//	   s3	    s4
//	conn2 <--> echo
//
//	   s5
//	proxy	from:s1,	to:s3
//		-- reads on s1	-- reads on s2
//
//	conn1:	s1.write()
//	echo:	s2.read()
//	echo:	s2.write()
//	proxy:	s1.read()
//	proxy:	s3.write()
//	proxy:	s3.read()
//	proxy:	s1.write()
//	echo:	s2.read()
//	echo:	s2.write()
//	proxy:	s1.read()
//	... loops until we cancel the transfer.
//
//	Note that we're looping only in one of the
//	sides (the one that initiated).
//
func TestProxyingInLoop(t *testing.T) {
	var msg = []byte("PING\r\n")
	var receiveBuffer = make([]byte, len(msg))
	var buf bytes.Buffer
	server := NewDumbTcpServer(&buf)
	defer server.Close()

	go func() {
		server.Listen()
	}()

	time.Sleep(200 * time.Millisecond)
	var addr = fmt.Sprintf("localhost:%d", server.GetPort())

	downstream, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	defer downstream.Close()

	upstream, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	defer upstream.Close()

	proxy, err := NewProxy(ProxyConfig{
		From: downstream,
		To:   upstream,
	})
	assert.NoError(t, err)

	go func() {
		assert.NoError(t, proxy.Transfer())
	}()

	time.Sleep(1000 * time.Millisecond)

	n, err := downstream.Write(msg)
	assert.NoError(t, err)
	assert.Equal(t, len(msg), n)

	n, err = upstream.Read(receiveBuffer)
	assert.NoError(t, err)
	assert.Equal(t, len(msg), n)
	assert.True(t, bytes.Count(buf.Bytes(), msg) > 4)
}
