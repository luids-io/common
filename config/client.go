// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"fmt"

	"github.com/luisguillenc/grpctls"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// ClientCfg stores grpc client preferences
type ClientCfg struct {
	RemoteURI string
	TLS       grpctls.ClientCfg
	Metrics   bool
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *ClientCfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	if short {
		pflag.StringVarP(&cfg.RemoteURI, aprefix+"uri", "r", cfg.RemoteURI, "URI to grpc service.")
	} else {
		pflag.StringVar(&cfg.RemoteURI, aprefix+"uri", cfg.RemoteURI, "URI to grpc service.")
	}
	pflag.StringVar(&cfg.TLS.CertFile, aprefix+"clientcert", cfg.TLS.CertFile, "Path to grpc client cert file.")
	pflag.StringVar(&cfg.TLS.KeyFile, aprefix+"clientkey", cfg.TLS.KeyFile, "Path to grpc client key file.")
	pflag.StringVar(&cfg.TLS.ServerCert, aprefix+"servercert", cfg.TLS.ServerCert, "Path to grpc server cert file.")
	pflag.StringVar(&cfg.TLS.ServerName, aprefix+"servername", cfg.TLS.ServerName, "Server name of grpc service for TLS check.")
	pflag.StringVar(&cfg.TLS.CACert, aprefix+"cacert", cfg.TLS.CACert, "Path to grpc CA cert file.")
	pflag.BoolVar(&cfg.TLS.UseSystemCAs, aprefix+"systemca", cfg.TLS.UseSystemCAs, "Use system CA pool for grpc check.")
	pflag.BoolVar(&cfg.Metrics, aprefix+"metrics", cfg.Metrics, "Enable metrics.")
}

// BindViper and bind to viper
func (cfg *ClientCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"uri")
	util.BindViper(v, aprefix+"clientcert")
	util.BindViper(v, aprefix+"clientkey")
	util.BindViper(v, aprefix+"servercert")
	util.BindViper(v, aprefix+"servername")
	util.BindViper(v, aprefix+"cacert")
	util.BindViper(v, aprefix+"systemca")
	util.BindViper(v, aprefix+"metrics")
}

// FromViper fill values from viper
func (cfg *ClientCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.RemoteURI = v.GetString(aprefix + "uri")
	cfg.TLS.CertFile = v.GetString(aprefix + "clientcert")
	cfg.TLS.KeyFile = v.GetString(aprefix + "clientkey")
	cfg.TLS.ServerCert = v.GetString(aprefix + "servercert")
	cfg.TLS.ServerName = v.GetString(aprefix + "servername")
	cfg.TLS.CACert = v.GetString(aprefix + "cacert")
	cfg.TLS.UseSystemCAs = v.GetBool(aprefix + "systemca")
	cfg.Metrics = v.GetBool(aprefix + "metrics")
}

// Empty returns true if configuration is empty
func (cfg ClientCfg) Empty() bool {
	if cfg.RemoteURI != "" {
		return false
	}
	if cfg.TLS.UseTLS() {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg ClientCfg) Validate() error {
	_, _, err := grpctls.ParseURI(cfg.RemoteURI)
	if err != nil {
		return err
	}
	if cfg.TLS.UseTLS() {
		return cfg.TLS.Validate()
	}
	return nil
}

// Dump configuration
func (cfg ClientCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
