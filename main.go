package main

import (
	"io"
	"log"
	"net"
)

type Proxy struct{}

func NewProxy() Proxy {
	return Proxy{}
}

func (p Proxy) Transfer(conn1, conn2 net.Conn) error {
	errChan := make(chan error, 1)
	go func() {
		_, err := io.Copy(conn2, conn1)
		conn1.Close()
		conn2.Close()
		errChan <- err
	}()

	_, err1 := io.Copy(conn1, conn2)
	conn1.Close()
	conn2.Close()
	err2 := <-errChan

	if err1 != nil {
		return err1
	}
	return err2
}

func handleConnection(conn net.Conn) {
	agent, err := net.Dial("tcp4", "0.0.0.0:8080")
	if err != nil {
		log.Panicf("couldn't dial port 8080", err)
	}

	proxy := NewProxy()
	proxy.Transfer(conn, agent)
}

func main() {
	ln, err := net.Listen("tcp4", ":8000")
	if err != nil {
		log.Panicf("couldn't listen on port 80", err)
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
