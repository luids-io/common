// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"errors"
	"fmt"
	"net"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/luids-io/common/config"
	"github.com/luids-io/core/grpctls"
	"github.com/luids-io/core/ipfilter"
)

// ErrURIServerExists defines error when a server for the uri was created.
var ErrURIServerExists = errors.New("uri server already exists")

// Server is a factory for a grpc server
func Server(cfg *config.ServerCfg) (net.Listener, *grpc.Server, error) {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	// check in server pool
	slis, srv, ok := serverPool.get(cfg.ListenURI)
	if ok {
		return slis, srv, ErrURIServerExists
	}
	// create server
	err := cfg.Validate()
	if err != nil {
		return nil, nil, fmt.Errorf("invalid server config: %v", err)
	}
	var creds credentials.TransportCredentials
	if cfg.TLS.UseTLS() {
		creds, err = grpctls.Creds(cfg.TLS)
		if err != nil {
			return nil, nil, fmt.Errorf("initializing TLS: %v", err)
		}
	}
	slis, err = grpctls.Listener(cfg.ListenURI)
	if err != nil {
		return nil, nil, fmt.Errorf("listening server: %v", err)
	}

	srv = grpc.NewServer(getGRPCServerOpts(
		creds, ipfilter.Whitelist(cfg.Allowed), cfg.Metrics)...)
	//write in pool
	serverPool.set(cfg.ListenURI, slis, srv)

	return slis, srv, nil
}

// setup grpc server middleware with server options
func getGRPCServerOpts(creds credentials.TransportCredentials, ipfilter ipfilter.Filter, metrics bool) []grpc.ServerOption {
	uinterceptors := make([]grpc.UnaryServerInterceptor, 0)
	sinterceptors := make([]grpc.StreamServerInterceptor, 0)
	if !ipfilter.Empty() {
		uinterceptors = append(uinterceptors, ipfilter.UnaryServerInterceptor)
		sinterceptors = append(sinterceptors, ipfilter.StreamServerInterceptor)
	}
	if metrics {
		uinterceptors = append(uinterceptors, grpc_prometheus.UnaryServerInterceptor)
		sinterceptors = append(sinterceptors, grpc_prometheus.StreamServerInterceptor)
	}
	//create options
	grpcopts := make([]grpc.ServerOption, 0)
	grpcopts = append(grpcopts,
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(uinterceptors...)),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(sinterceptors...)))
	if creds != nil {
		grpcopts = append(grpcopts, grpc.Creds(creds))
	}
	return grpcopts
}

type grpcItem struct {
	listener net.Listener
	server   *grpc.Server
}

type grpcPool struct {
	items map[string]grpcItem
}

func (p *grpcPool) get(uri string) (net.Listener, *grpc.Server, bool) {
	item, ok := p.items[uri]
	if !ok {
		return nil, nil, false
	}
	return item.listener, item.server, true
}

func (p *grpcPool) set(uri string, lis net.Listener, srv *grpc.Server) {
	p.items[uri] = grpcItem{listener: lis, server: srv}
}

var serverMutex sync.Mutex
var serverPool grpcPool

func init() {
	serverPool = grpcPool{items: make(map[string]grpcItem)}
}
