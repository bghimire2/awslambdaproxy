package server

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/ginuerzh/gost"
)

const (
	forwardProxy = "localhost:8082"
)

// LocalProxy is proxy listener and where to forward
type LocalProxy struct {
	listeners    []string
	forwardProxy string
}

func (l *LocalProxy) run() {
	baseCfg := &baseConfig{}
	baseCfg.route.ChainNodes = []string{l.forwardProxy}
	baseCfg.route.ServeNodes = l.listeners

	cert, err := gost.GenCertificate()
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	gost.DefaultTLSConfig = tlsConfig

	var routers []router
	rts, err := baseCfg.route.GenRouters()
	if err != nil {
		log.Fatal(err)
	}
	routers = append(routers, rts...)

	if len(routers) == 0 {
		log.Fatalln("invalid config", err)
	}
	for i := range routers {
		go routers[i].Serve()
	}
}

// NewLocalProxy starts a local proxy that will forward to proxy running in Lambda
func NewLocalProxy(listeners []string, debugProxy bool, bypass string) (*LocalProxy, error) {
	if debugProxy {
		gost.SetLogger(&gost.LogLogger{})
	}
	fproxy := forwardProxy
	if bypass != "" {
		fproxy += fmt.Sprintf("?bypass=%v", bypass)
	}
	l := &LocalProxy{
		listeners:    listeners,
		forwardProxy: fproxy,
	}
	go l.run()
	return l, nil
}
