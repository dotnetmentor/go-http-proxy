# go-http-proxy

Proxy HTTP requests to a specified host target. Handles Location response header re-writes and basic auth for protecting the target host.

## Usage

### Build from source

```bash
git clone git@github.com:dotnetmentor/go-http-proxy.git
cd ./go-http-proxy
go build -o go-http-proxy
go-http-proxy -addr=localhost:8080 -host=localhost:2113 -rewrite-location-header -basic=admin:changeit -v
```

### Docker

```bash
docker build -t dotnetmentor/go-http-proxy:latest .
docker run --rm -p 8080:8080 dotnetmentor/go-http-proxy:latest -host=<hostname>:<port> -v
```
