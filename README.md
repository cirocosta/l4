# l4

> TCP load-balancer that collects all sorts of metrics

## Features

- Graceful shutdown of established connections
- Established connections count 
- Connection mapping with name resolution
- Tx and Rx stats

## Overview

Configuration:

```
# config.yml
port: 80
servers:
- address: 'http://192.168.0.103:8081'
- address: '192.168.0.103:8082'         
- address: 'http://nginx'                 
```

Once a connection arrives to `l4`, picks one from the list.

