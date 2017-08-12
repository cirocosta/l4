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

// down --> proxy
// proxy --> upstream
func TestProxying(t *testing.T) {
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

	time.Sleep(200 * time.Millisecond)

	n, err := downstream.Write(msg)
	assert.NoError(t, err)
	assert.Equal(t, len(msg), n)

	n, err = upstream.Read(receiveBuffer)
	assert.NoError(t, err)
	assert.Equal(t, len(msg), n)
	assert.True(t, bytes.Count(buf.Bytes(), msg) > 50)
}
