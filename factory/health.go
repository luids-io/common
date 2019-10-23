// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"
	"net"

	"github.com/luisguillenc/httphealth"
	"github.com/luisguillenc/ipfilter"
	"github.com/luisguillenc/yalogi"

	"github.com/luids-io/common/config"
	"github.com/luids-io/common/util"
)

// Health is a factory for an http server
func Health(cfg *config.HealthCfg, srv httphealth.Pingable, logger yalogi.Logger) (net.Listener, *httphealth.Server, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, nil, fmt.Errorf("invalid health config: %v", err)
	}
	hlis, err := util.Listener(cfg.ListenURI)
	if err != nil {
		return nil, nil, fmt.Errorf("listening health: %v", err)
	}
	health := httphealth.New(srv,
		httphealth.SetLogger(logger),
		httphealth.Metrics(cfg.Metrics),
		httphealth.SetIPFilter(ipfilter.Whitelist(cfg.Allowed)))
	return hlis, health, nil
}