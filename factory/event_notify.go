// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"
	"time"

	"github.com/luids-io/api/event"
	"github.com/luids-io/api/event/notifybuffer"
	"github.com/luids-io/common/config"
	"github.com/luids-io/core/apiservice"
	"github.com/luids-io/core/yalogi"
)

// EventNotify is a factory for an event notifier.
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

// EventNotifyBuffer is a factory for an event buffer.
func EventNotifyBuffer(cfg *config.EventNotifyCfg, registry apiservice.Discover, logger yalogi.Logger) (event.NotifyBuffer, error) {
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
	var output event.NotifyBuffer
	output = notifybuffer.Notifier(client, logger)
	if cfg.WaitDups > 0 && cfg.Buffer > 0 {
		output = notifybuffer.NewWaitDups(output, cfg.Buffer, time.Duration(cfg.WaitDups)*time.Millisecond)
	}
	if cfg.Buffer > 0 {
		output = notifybuffer.NewQueue(output, cfg.Buffer)
	}
	return output, nil
}
