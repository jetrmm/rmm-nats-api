nats-api
===

This is the NATS API server for JetRMM. It runs alongside the server to handle agent RPC communications.

## Building the NATS API server

```shell
env CGO_ENABLED=0 GOARCH=amd64 go build -ldflags "-s -w" -o out\nats-api
```
