// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"errors"
	"fmt"
	"net"

	"github.com/luids-io/common/util"
	"github.com/luisguillenc/grpctls"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ServerCfg stores server preferences
type ServerCfg struct {
	ListenURI string
	Allowed   []string
	TLS       grpctls.ServerCfg
	Metrics   bool
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *ServerCfg) SetPFlags(prefix string) {
	short := true
	aprefix := ""
	if prefix != "" {
		short = false
		aprefix = prefix + "."
	}
	if short {
		pflag.StringVarP(&cfg.ListenURI, aprefix+"listenuri", "l", cfg.ListenURI, "Server socket.")
	} else {
		pflag.StringVar(&cfg.ListenURI, aprefix+"listenuri", cfg.ListenURI, "Server socket.")
	}
	pflag.StringSliceVar(&cfg.Allowed, aprefix+"allowed", cfg.Allowed, "List of allowed IPs or CIDRs.")
	pflag.StringVar(&cfg.TLS.CertFile, aprefix+"certfile", cfg.TLS.CertFile, "Path to server cert file.")
	pflag.StringVar(&cfg.TLS.KeyFile, aprefix+"keyfile", cfg.TLS.KeyFile, "Path to server key file.")
	pflag.StringVar(&cfg.TLS.CACert, aprefix+"cacert", cfg.TLS.CACert, "Path to CA cert file.")
	pflag.BoolVar(&cfg.TLS.ClientAuth, aprefix+"clientauth", cfg.TLS.ClientAuth, "Require client auth.")
	pflag.BoolVar(&cfg.Metrics, aprefix+"metrics", cfg.Metrics, "Enable metrics.")
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *ServerCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"listenuri")
	util.BindViper(v, aprefix+"allowed")
	util.BindViper(v, aprefix+"certfile")
	util.BindViper(v, aprefix+"keyfile")
	util.BindViper(v, aprefix+"cacert")
	util.BindViper(v, aprefix+"clientauth")
	util.BindViper(v, aprefix+"metrics")
}

// FromViper fill values from viper
func (cfg *ServerCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.ListenURI = v.GetString(aprefix + "listenuri")
	cfg.Allowed = v.GetStringSlice(aprefix + "allowed")
	cfg.TLS.CertFile = v.GetString(aprefix + "certfile")
	cfg.TLS.KeyFile = v.GetString(aprefix + "keyfile")
	cfg.TLS.CACert = v.GetString(aprefix + "cacert")
	cfg.TLS.ClientAuth = v.GetBool(aprefix + "clientauth")
	cfg.Metrics = v.GetBool(aprefix + "metrics")
}

// Empty returns true if configuration is empty
func (cfg ServerCfg) Empty() bool {
	if cfg.ListenURI != "" {
		return false
	}
	if cfg.Allowed != nil && len(cfg.Allowed) > 0 {
		return false
	}
	if cfg.TLS.UseTLS() {
		return false
	}
	if cfg.Metrics {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg *ServerCfg) Validate() error {
	if cfg.ListenURI == "" {
		return errors.New("listen uri is required")
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
	if cfg.TLS.UseTLS() {
		return cfg.TLS.Validate()
	}
	return nil
}

// Dump configuration
func (cfg ServerCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
