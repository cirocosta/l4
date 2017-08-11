package lib

import (
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
