// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"errors"
	"fmt"
	"net"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// HealthCfg stores http health server preferences
type HealthCfg struct {
	ListenURI string
	Allowed   []string
	Metrics   bool
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *HealthCfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	pflag.StringVar(&cfg.ListenURI, aprefix+"listenuri", cfg.ListenURI, "Health and metrics socket.")
	pflag.BoolVar(&cfg.Metrics, aprefix+"metrics", cfg.Metrics, "Expose prometheus metrics.")
	pflag.StringSliceVar(&cfg.Allowed, aprefix+"allowed", cfg.Allowed, "List of allowed IPs or CIDRs.")
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *HealthCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"listenuri")
	util.BindViper(v, aprefix+"metrics")
	util.BindViper(v, aprefix+"allowed")
}

// FromViper fill values from viper
func (cfg *HealthCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.ListenURI = v.GetString(aprefix + "listenuri")
	cfg.Metrics = v.GetBool(aprefix + "metrics")
	cfg.Allowed = v.GetStringSlice(aprefix + "allowed")
}

// Empty returns true if configuration is empty
func (cfg HealthCfg) Empty() bool {
	if cfg.ListenURI != "" {
		return false
	}
	if cfg.Allowed != nil && len(cfg.Allowed) > 0 {
		return false
	}
	if cfg.Metrics {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg HealthCfg) Validate() error {
	if cfg.ListenURI == "" {
		return errors.New("listenuri is required")
	}
	_, _, err := util.ParseListenURI(cfg.ListenURI)
	if err != nil {
		return err
	}
	for _, item := range cfg.Allowed {
		_, _, err = net.ParseCIDR(item)
		if err != nil {
			ip := net.ParseIP(item)
			if ip == nil {
				return fmt.Errorf("value '%v' is not a valid ip or cidr", item)
			}
		}
	}
	return nil
}

// Dump configuration
func (cfg HealthCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
