// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"

	"github.com/luids-io/api/event"
	"github.com/luids-io/common/config"
	"github.com/luids-io/core/apiservice"
)

// EventNotify is a factory for an event notifier using apiservice
func EventNotify(cfg *config.EventNotifyCfg, registry apiservice.Discover) (event.Notifier, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid event notify config: %v", err)
	}
	svc, ok := registry.GetService(cfg.Service)
	if !ok {
		return nil, fmt.Errorf("can't find service '%s'", cfg.Service)
	}
	client, ok := svc.(event.Notifier)
	if !ok {
		return nil, fmt.Errorf("service '%s' is not a notifier", cfg.Service)
	}
	return client, nil
}
