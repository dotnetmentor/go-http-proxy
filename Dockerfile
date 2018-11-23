FROM golang:alpine as builder
RUN apk add --no-cache git
COPY . $GOPATH/src/go-http-proxy/
WORKDIR $GOPATH/src/go-http-proxy/
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/go-http-proxy

FROM scratch
COPY --from=builder /go/bin/go-http-proxy /go/bin/go-http-proxy
EXPOSE 8080
ENTRYPOINT ["/go/bin/go-http-proxy"]
