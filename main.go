package main

import (
	"log"
	"net"
	"time"

	. "github.com/cirocosta/l4/lib"
)

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

func main() {
	ln, err := net.Listen("tcp4", ":8000")
	if err != nil {
		log.Panicf("couldn't listen on port 8000", err)
	}

	ln = NewGracefulListener(GracefulListenerConfig{
		Listener: ln,
	})

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
