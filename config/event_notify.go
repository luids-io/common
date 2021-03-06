// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// EventNotifyCfg stores event-notify client configuration.
type EventNotifyCfg struct {
	Service  string
	Instance string
	Buffer   int
	WaitDups int
}

// SetPFlags setups posix flags for commandline configuration.
func (cfg *EventNotifyCfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	pflag.StringVar(&cfg.Service, aprefix+"service", cfg.Service, "API Service ID.")
	pflag.StringVar(&cfg.Instance, aprefix+"instance", cfg.Instance, "Instance name.")
	pflag.IntVar(&cfg.Buffer, aprefix+"buffer", cfg.Buffer, "Buffer size.")
	pflag.IntVar(&cfg.Buffer, aprefix+"waitdups", cfg.Buffer, "Wait for duplicates (in milliseconds).")
}

// BindViper setups posix flags for commandline configuration and bind to viper.
func (cfg *EventNotifyCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"service")
	util.BindViper(v, aprefix+"instance")
	util.BindViper(v, aprefix+"buffer")
	util.BindViper(v, aprefix+"waitdups")
}

// FromViper fill values from viper.
func (cfg *EventNotifyCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.Service = v.GetString(aprefix + "service")
	cfg.Instance = v.GetString(aprefix + "instance")
	cfg.Buffer = v.GetInt(aprefix + "buffer")
	cfg.WaitDups = v.GetInt(aprefix + "waitdups")
}

// Empty returns true if configuration is empty.
func (cfg EventNotifyCfg) Empty() bool {
	if cfg.Service != "" {
		return false
	}
	return true
}

// Validate checks that configuration is ok.
func (cfg EventNotifyCfg) Validate() error {
	if cfg.Service == "" {
		return errors.New("service name required")
	}
	if cfg.Buffer <= 0 {
		return errors.New("invalid buffer size")
	}
	if cfg.WaitDups < 0 {
		return errors.New("invalid waitdups")
	}
	return nil
}

// Dump configuration.
func (cfg EventNotifyCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
