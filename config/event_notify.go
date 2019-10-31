// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// EventNotifyCfg stores event notify client (using apiservice)
type EventNotifyCfg struct {
	Service string
	Buffer  int
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *EventNotifyCfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	pflag.StringVar(&cfg.Service, aprefix+"service", cfg.Service, "API Service ID.")
	pflag.IntVar(&cfg.Buffer, aprefix+"buffer", cfg.Buffer, "Buffer size.")
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *EventNotifyCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"service")
	util.BindViper(v, aprefix+"buffer")
}

// FromViper fill values from viper
func (cfg *EventNotifyCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.Service = v.GetString(aprefix + "service")
	cfg.Buffer = v.GetInt(aprefix + "buffer")
}

// Empty returns true if configuration is empty
func (cfg EventNotifyCfg) Empty() bool {
	if cfg.Service != "" {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg EventNotifyCfg) Validate() error {
	if cfg.Service == "" {
		return errors.New("service name required")
	}
	if cfg.Buffer <= 0 {
		return errors.New("invalid buffer size")
	}
	return nil
}

// Dump configuration
func (cfg EventNotifyCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
