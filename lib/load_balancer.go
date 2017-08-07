package lib

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type server struct {
	address          string
	activeProxies    []Proxy
	totalConnections uint64
	totalRx          uint64
	totalTx          uint64
}

type LoadBalancer struct {
	servers []*server
	nextIdx uint64
	port    int
}

type LoadBalancerConfig struct {
	Servers []string
	Port    int
}

func NewLoadBalancer(cfg LoadBalancerConfig) (lb LoadBalancer, err error) {
	if len(cfg.Servers) == 0 {
		err = errors.Errorf("At least one server must be specified")
		return
	}

	return
}

func (lb *LoadBalancer) Load(addresses []string) (err error) {
	if len(addresses) == 0 {
		err = errors.Errorf("'Load' must receive at least one server")
		return
	}

	var servers = make([]*server, len(addresses))
	for ndx, address := range addresses {
		serversList[ndx] = &server{
			address: address,
		}
	}

	lb.servers = servers
}

func (lb *LoadBalancer) handle(conn net.Conn) {
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

func (lb *LoadBalancer) Listen() (err error) {
	ln, err := net.Listen("tcp4", fmt.Sprintf(":%d", lb.port))
	if err != nil {
		err = errors.Wrapf(err, "couldn't listen on port %d", lb.port)
		return
	}

	ln = NewGracefulListener(GracefulListenerConfig{
		Listener: ln,
	})

	for {
		conn, err := ln.Accept()
		if err != nil {
			// ops
		}

		lb.handle(conn)
	}

}

func (lb *LoadBalancer) Stop() (err error) {
	// for each proxy, make it stop.
}
