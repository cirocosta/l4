package main

import (
	"log"
	"net"
	"time"

	"github.com/alexflint/go-arg"

	. "github.com/cirocosta/l4/lib"
)

type config struct {
	Port    int      `arg:"-p,help:port to listen to"`
	Config  string   `arg:"-c,help:configuration file to use"`
	Servers []string `arg:"positional"`
}

var (
	args = &config{Port: 80}
)

func main() {
	lb, err := NewLoadBalancer(LoadBalancerConfig{
		Servers: []string{
			"127.0.0.1:8080",
		},
	})
	if err != nil {
		log.Panicf("ERROR: Couldn't listen for connections %+v\n", err)
	}
}
