package main

import (
	"net"
	"time"
)

type IoStats struct {
	Tx uint64
	Rx uint64
}

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

type GracefulConn struct {
	net.Conn
	ln *GracefulListener
}
