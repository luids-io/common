// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// LoggerCfg stores logger configuration preferences
type LoggerCfg struct {
	Level  string
	Format string
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *LoggerCfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	pflag.StringVar(&cfg.Level, aprefix+"level", cfg.Level, "Log level.")
	pflag.StringVar(&cfg.Format, aprefix+"format", cfg.Format, "Log format.")
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *LoggerCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"level")
	util.BindViper(v, aprefix+"fomat")
}

// FromViper fill values from viper
func (cfg *LoggerCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.Level = v.GetString(aprefix + "level")
	cfg.Format = v.GetString(aprefix + "format")
}

// Empty returns true if configuration is empty
func (cfg LoggerCfg) Empty() bool {
	if cfg.Level != "" {
		return false
	}
	if cfg.Format != "" {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg LoggerCfg) Validate() error {
	switch strings.ToLower(cfg.Format) {
	case "": //ok
	case "json": //ok
	case "text": //ok
	case "log": //ok
	default:
		return errors.New("invalid format value")
	}
	switch strings.ToLower(cfg.Level) {
	case "error":
		return nil
	case "warn", "warning":
		return nil
	case "info":
		return nil
	case "debug":
		return nil
	}
	return errors.New("invalid level value")
}

// Dump configuration
func (cfg LoggerCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
