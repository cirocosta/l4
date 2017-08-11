# l4

> TCP load-balancer that collects all sorts of metrics

## Features

- Graceful shutdown of established connections
- Established connections count 
- Connection mapping with name resolution
- Tx and Rx stats


## Overview

### CLI

```
Usage: l4 [--port PORT] [--config CONFIG] [--debug] [SERVERS [SERVERS ...]]

Positional arguments:
  SERVERS

Options:
  --port PORT, -p PORT   port to listen to [default: 3000]
  --config CONFIG, -c CONFIG
                         configuration file to use
  --debug, -d            enables debug mode
  --help, -h             display this help and exit

Example:
  l4 \
	--debug \
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

