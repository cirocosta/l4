<h1 align="center">l4 📂  </h1>

<h5 align="center">Minimal TCP load-balancer</h5>

<br/>

[![Build Status](https://travis-ci.org/cirocosta/l4.svg?branch=master)](https://travis-ci.org/cirocosta/l4)


## Features

- Graceful shutdown of established connections
- Established connections count 
- Connection mapping with name resolution
- Tx and Rx stats
- SNI


## Overview

### CLI

```
Usage: l4 [opts] SERVERS [SERVERS ...]

Positional arguments:
  SERVERS                list of <ip>:<port> or <domain>:<port> servers to connect to

Options:
  --debug                enables debug mode
  --port PORT            port to listen to [default: 1337]
  --server-name SERVER-NAME
                         server name to use on tls connections
  --tls-connect          connects to backends via TLS
  --tls-key-log TLS-KEY-LOG
                         log tls master secrets to file
  --tls-listen           listens for TLS connections
  --tls-listen-cert TLS-LISTEN-CERT
                         certificate to use when listening to TLS conns
  --tls-listen-key TLS-LISTEN-KEY
                         private key to use when listening to TLS conns
  --tls-skip-verification
                         skips certificate verification
  --help, -h             display this help and exit

Example:
  l4 \
	--port 3000 \
	127.0.0.1:3000 \
	127.0.0.1:3001
```

### Docker

To run `l4` as a docker container all you need to do is use `cirocosta/l4` and specify the same parameters that are used in the CLI.

Example:

```
# create a network 
docker network create mynet

# create two nginx containers in this network
# without exposed ports
docker run --detach --network mynet --name nginx1 nginx:alpine
docker run --detach --network mynet --name nginx2 nginx:alpine

# Create the l4 load-balancer specifying the name of
# each of the containers
docker run \
	--detach \
	--publish 80:3000 \
	cirocosta/l4 \
	nginx1:80 nginx2:80
```

