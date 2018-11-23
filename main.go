package main

import (
	"crypto/subtle"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/elazarl/goproxy"
)

const realm = "go-proxy"
const addressFlagKey = "addr"
const hostFlagKey = "host"

func main() {
	verbose := flag.Bool("v", false, "log proxy request to stdout")
	addr := flag.String(addressFlagKey, ":8080", "proxy listen address")
	host := flag.String(hostFlagKey, "", "proxy requests to host address (example: \"localhost:2113\")")
	rewriteLocationHeader := flag.Bool("rewrite-location-header", false, "allow modifying the Location response header")
	basicAuth := flag.String("basic", "", "protect using basic auth (example: \"username:password\")")
	flag.Parse()

	if *addr == "" {
		flag.Usage()
		log.Fatalln(fmt.Sprintf("-%s must be specified", addressFlagKey))
	}

	if *host == "" {
		flag.Usage()
		log.Fatalln(fmt.Sprintf("-%s must be specified", hostFlagKey))
	}

	username := ""
	password := ""

	if *basicAuth != "" && strings.Contains(*basicAuth, ":") {
		parts := strings.Split(*basicAuth, ":")
		if len(parts[0]) >= 1 && len(parts[1]) >= 1 {
			username = parts[0]
			password = parts[1]
		}
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if username != "" && password != "" {
			user, pass, ok := req.BasicAuth()
			if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
				res := goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusUnauthorized, "Don't waste your time!")
				res.Header.Set("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", realm))
				return req, res
			}
		}
		return req, nil
	})

	proxy.OnResponse().DoFunc(func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		cl := res.Header.Get("Location")
		if cl != "" && *rewriteLocationHeader {
			a := *addr
			if strings.HasPrefix(a, ":") {
				a = "127.0.0.1" + a
			}
			u, _ := url.Parse(cl)
			nl := strings.Replace(cl, u.Host, a, -1)

			if nl != cl {
				ctx.Logf("Modifying HTTP response header (Location: %s -> %s)", cl, nl)
				res.Header.Set("Location", nl)
			}
		}
		return res
	})

	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = *host
		proxy.ServeHTTP(w, req)
	})

	proxy.Logger.Println(fmt.Sprintf("Started proxy server on \"%s\"...", *addr))
	proxy.Logger.Println(fmt.Sprintf("Proxying requests to \"%s\"", *host))

	log.Fatal(http.ListenAndServe(*addr, proxy))
}
