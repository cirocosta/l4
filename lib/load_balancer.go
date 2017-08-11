package lib

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
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
	nextIdx int
	port    int
	logger  zerolog.Logger
}

type LoadBalancerConfig struct {
	Port  int
	Debug bool
}

func NewLoadBalancer(cfg LoadBalancerConfig) (lb LoadBalancer, err error) {
	if cfg.Port == 0 {
		err = errors.Errorf("a port != 0 must be specified")
		return
	}

	if cfg.Debug {
		lb.logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		lb.logger = zerolog.New(os.Stderr)
	}

	lb.port = cfg.Port
	return
}

func (lb *LoadBalancer) Load(addresses []string) (err error) {
	lb.logger.Info().
		Int("n-addresses", len(addresses)).
		Msg("loading configuration")

	var servers []*server

	if len(addresses) == 0 {
		err = errors.Errorf("must specify at least one server")
		return
	}

	servers = make([]*server, len(addresses))
	for ndx, address := range addresses {
		servers[ndx] = &server{
			address: address,
		}
	}

	lb.servers = servers
	return
}

func (lb *LoadBalancer) Listen() (err error) {
	lb.logger.Info().
		Int("port", lb.port).
		Msg("listening")

	ln, err := net.Listen("tcp4", fmt.Sprintf(":%d", lb.port))
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't listen on port %d", lb.port)
		return
	}

	ln = NewGracefulListener(GracefulListenerConfig{
		Listener: ln,
	})

	for {
		conn, err := ln.Accept()
		if err != nil {
			lb.logger.Error().
				Err(err).
				Msg("errored accepting connection")
			continue
		}

		go lb.handle(conn, lb.servers[lb.nextIdx%len(lb.servers)])
	}
}

func (lb *LoadBalancer) handle(conn net.Conn, s *server) {
	var logger = lb.logger.With().
		Str("local", conn.LocalAddr().String()).
		Str("upstream", s.address).
		Str("id", xid.New().String()).
		Logger()

	logger.Info().Msg("dialing")
	agent, err := net.Dial("tcp4", s.address)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("couldn't dial server")
		return
	}

	proxy, err := NewProxy(ProxyConfig{
		To:                agent,
		From:              conn,
		ConnectionTimeout: 10 * time.Second,
	})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("couldn't create proxy")
		return
	}

	logger.Info().Msg("proxying")
	err = proxy.Transfer()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("errored transferring between connections")
		return
	}

	logger.Info().Msg("finished")
}

func (lb *LoadBalancer) Stop() (err error) {
	return
}
