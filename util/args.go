// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package util

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// BindViper binds pflagKey to viper
func BindViper(v *viper.Viper, pflagKey string) {
	v.BindPFlag(pflagKey, pflag.Lookup(pflagKey))
}

// IsValid returns true if value is into the valid set
func IsValid(value string, valid []string) bool {
	for _, s := range valid {
		if value == s {
			return true
		}
	}
	return false
}
