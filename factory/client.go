// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"

	"github.com/luids-io/common/config"
	"github.com/luids-io/core/utils/grpctls"
)

// ClientConn is a factory for a grpc dial collector
func ClientConn(cfg *config.ClientCfg) (*grpc.ClientConn, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid client config: %v", err)
	}
	opts := make([]grpc.DialOption, 0)
	if cfg.Metrics {
		opts = append(opts, grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor))
		opts = append(opts, grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor))
	}
	//create dial
	dial, err := grpctls.Dial(cfg.RemoteURI, cfg.TLS, opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot dial with %s: %v", cfg.RemoteURI, err)
	}
	return dial, err
}
