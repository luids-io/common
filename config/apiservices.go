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

// APIServicesCfg stores services
type APIServicesCfg struct {
	ConfigDirs  []string
	ConfigFiles []string
	CertsDir    string
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *APIServicesCfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	pflag.StringSliceVar(&cfg.ConfigDirs, aprefix+"dirs", cfg.ConfigDirs, "Configuration dirs.")
	pflag.StringSliceVar(&cfg.ConfigFiles, aprefix+"files", cfg.ConfigFiles, "Configuration files.")
	pflag.StringVar(&cfg.CertsDir, aprefix+"certsdir", cfg.CertsDir, "Base path to certificate files.")
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *APIServicesCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"dirs")
	util.BindViper(v, aprefix+"files")
	util.BindViper(v, aprefix+"certsdir")
}

// FromViper fill values from viper
func (cfg *APIServicesCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.ConfigDirs = v.GetStringSlice(aprefix + "dirs")
	cfg.ConfigFiles = v.GetStringSlice(aprefix + "files")
	cfg.CertsDir = v.GetString(aprefix + "certsdir")
}

// Empty returns true if configuration is empty
func (cfg APIServicesCfg) Empty() bool {
	if len(cfg.ConfigFiles) > 0 {
		return false
	}
	if len(cfg.ConfigDirs) > 0 {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg APIServicesCfg) Validate() error {
	empty := true
	for _, file := range cfg.ConfigFiles {
		if !util.FileExists(file) {
			return fmt.Errorf("config file '%s' doesn't exists", file)
		}
		if !strings.HasSuffix(file, ".json") {
			return fmt.Errorf("config file '%s' without .json extension", file)
		}
		empty = false
	}
	for _, dir := range cfg.ConfigDirs {
		if !util.DirExists(dir) {
			return fmt.Errorf("config dir '%s' doesn't exists", dir)
		}
		empty = false
	}
	if empty {
		return errors.New("config required")
	}

	if cfg.CertsDir != "" {
		if !util.DirExists(cfg.CertsDir) {
			return fmt.Errorf("certificates dir '%v' doesn't exists", cfg.CertsDir)
		}
	}
	return nil
}

// Dump configuration
func (cfg APIServicesCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
