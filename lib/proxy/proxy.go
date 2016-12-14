package proxy

import (
	"errors"
	"fmt"
	"github.com/elazarl/goproxy"
	"net/http"
	// "regexp"
	"sync"
)

type Proxy struct {
	Proxy   *goproxy.ProxyHttpServer
	Addr    string
	Port    string
	Enabled bool
}

func (p *Proxy) start() error {
	p.Proxy.OnRequest().DoFunc(p.handleRequest)
	p.Proxy.OnResponse().DoFunc(p.handleResponse)
	http.ListenAndServe(p.Addr+":"+p.Port, p.Proxy)
	//It should never return
	return nil
}

func (p *Proxy) Run(wg sync.WaitGroup) error {
	if p.Port == "" {
		return errors.New("proxy: Port must be specified")
	}
	go p.start()
	wg.Done()
	return nil
}

// Proxy.handleRequest is called every time the proxy gets an HTTP request from the client
func (p *Proxy) handleRequest(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if !p.Enabled {
		return r, nil
	}
	fmt.Printf("Method: %s %s\nProtocol: %s\nHost: %s\n", r.Method, r.URL, r.Proto, r.Host)
	fmt.Printf("%s", r.Body)
	return r, nil
}

// Proxy.handleResponse is called every time the proxy gets an HTTP response
func (p *Proxy) handleResponse(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if !p.Enabled {
		return r
	}
	return r
}

// NewIntercepterProxy creates a Proxy struct listening to the adress and port passed to the function
func NewIntercepterProxy(addr string, port string) *Proxy {
	ret := &Proxy{}
	ret.Proxy = goproxy.NewProxyHttpServer()
	ret.Addr = addr
	ret.Port = port
	ret.Enabled = false

	return ret
}
