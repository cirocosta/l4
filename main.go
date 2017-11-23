package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"

	. "github.com/cirocosta/l4/lib"
)

type config struct {
	Debug               bool     `arg:"--debug,env,help:enables debug mode"`
	Port                int      `arg:"--port,env,help:port to listen to"`
	ServerName          string   `arg:"--server-name,env,help:server name to use on tls connections"`
	Servers             []string `arg:"required,positional,help:list of <ip>:<port> or <domain>:<port> servers to connect to"`
	TlsConnect          bool     `arg:"--tls-connect,env,help:connects to backends via TLS"`
	TlsKeyLog           string   `arg:"--tls-key-log,env:SSLKEYLOGFILE,help:log tls master secrets to file"`
	TlsListen           bool     `arg:"--tls-listen,env,help:listens for TLS connections"`
	TlsListenCert       string   `arg:"--tls-listen-cert,help:certificate to use when listening to TLS conns"`
	TlsListenKey        string   `arg:"--tls-listen-key,help:private key to use when listening to TLS conns"`
	TlsSkipVerification bool     `arg:"--tls-skip-verification,env:SKIP_VERIFICATION,help:skips certificate verification"`
}

var (
	args = &config{
		Debug:               false,
		Port:                1337,
		ServerName:          "",
		TlsConnect:          false,
		TlsKeyLog:           "",
		TlsListen:           false,
		TlsListenCert:       "",
		TlsListenKey:        "",
		TlsSkipVerification: false,
	}
	argParser *arg.Parser
	err       error
)

func abort(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "ERROR: \n%+v\n", err)
	os.Exit(1)
}

func main() {
	arg.MustParse(args)

	lb, err := NewLoadBalancer(LoadBalancerConfig{
		Port:  args.Port,
		Debug: args.Debug,
	})
	abort(err)

	err = lb.Load(args.Servers)
	abort(err)

	err = lb.Listen()
	abort(err)
}
