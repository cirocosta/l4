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
port: 8080
servers:
  - '127.0.0.1:8001'
  - '127.0.0.1:8002'
  - '127.0.0.1:8003'
```

Once a connection arrives to `l4`, picks one from the list.

