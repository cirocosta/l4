package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"

	. "github.com/cirocosta/l4/lib"
)

type config struct {
	Port    int      `arg:"-p,env,help:port to listen to"`
	Config  string   `arg:"-c,env,help:configuration file to use"`
	Debug   bool     `arg:"-d,env,help:enables debug mode"`
	Servers []string `arg:"positional"`
}

var (
	args      = &config{Port: 3000}
	argParser *arg.Parser
	err       error
)

func main() {
	argParser = arg.MustParse(args)
	if len(args.Servers) == 0 {
		argParser.Fail("At least one server must be specified (positional argument).")
	}

	lb, err := NewLoadBalancer(LoadBalancerConfig{
		Port:  args.Port,
		Debug: args.Debug,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Couldn't instantiate load-balancer.\n"+
			"%+v\n", err)
		os.Exit(1)
	}

	err = lb.Load(args.Servers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Couldn't make load-balancer load server configuration.\n"+
			"%+v\n", err)
		os.Exit(1)
	}

	err = lb.Listen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed listening on port %d.\n"+
			"%+v\n", args.Port, err)
		os.Exit(1)
	}
}
