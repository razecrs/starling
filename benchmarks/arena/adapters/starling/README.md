# Starling arena adapter

This adapter uses `Client.HandleGatewayFrame`, the same typed dispatcher and
internal state consumers used by the live gateway. It never opens a network
connection.

```sh
bash benchmarks/arena/adapters/starling/run.sh prepare
bash benchmarks/arena/adapters/starling/run.sh verify
bash benchmarks/arena/adapters/starling/run.sh bench message_handled 10000 5 30
```

Supported: all six canonical correctness/benchmark workloads. The adapter does
not edit Starling or embed fixture copies.
