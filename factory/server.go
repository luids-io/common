// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/luisguillenc/grpctls"
	"github.com/luisguillenc/ipfilter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/luids-io/common/config"
)

// Server is a factory for a grpc server
func Server(cfg *config.ServerCfg) (net.Listener, *grpc.Server, error) {
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
	slis, err := grpctls.Listener(cfg.ListenURI)
	if err != nil {
		return nil, nil, fmt.Errorf("listening server: %v", err)
	}

	srv := grpc.NewServer(getGRPCServerOpts(
		creds, ipfilter.Whitelist(cfg.Allowed), cfg.Metrics)...)

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
