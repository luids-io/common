// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"

	"github.com/luids-io/common/config"
	"github.com/luids-io/common/util"
	"github.com/luids-io/core/apiservice"
	"github.com/luids-io/core/yalogi"
)

// APIServices is a factory of an APIService Registry.
func APIServices(cfg *config.APIServicesCfg, logger yalogi.Logger) (*apiservice.Registry, error) {
	if cfg.Empty() {
		return apiservice.NewRegistry(), nil
	}
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	defs, err := getServiceDefs(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("loading servicedefs: %v", err)
	}
	services := apiservice.NewRegistry()
	for _, def := range defs {
		if def.Disabled {
			logger.Debugf("'%s' is disabled", def.ID)
			continue
		}
		svc, err := apiservice.Build(def, logger)
		if err != nil {
			return nil, fmt.Errorf("building '%s': %v", def.ID, err)
		}
		services.Register(def.ID, svc)
	}
	return services, nil
}

// APIAutoloader is a factory of an APIService Autoloader .
func APIAutoloader(cfg *config.APIServicesCfg, logger yalogi.Logger) (*apiservice.Autoloader, error) {
	if cfg.Empty() {
		return apiservice.NewAutoloader([]apiservice.ServiceDef{}, apiservice.SetLogger(logger)), nil
	}
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	defs, err := getServiceDefs(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("loading servicedefs: %v", err)
	}
	autodefs := make([]apiservice.ServiceDef, 0, len(defs))
	for _, def := range defs {
		if def.Disabled {
			logger.Debugf("'%s' is disabled", def.ID)
			continue
		}
		autodefs = append(autodefs, def)
	}
	return apiservice.NewAutoloader(autodefs, apiservice.SetLogger(logger)), nil
}

func getServiceDefs(cfg *config.APIServicesCfg, logger yalogi.Logger) ([]apiservice.ServiceDef, error) {
	dbFiles, err := util.GetFilesDB("json", cfg.ConfigFiles, cfg.ConfigDirs)
	if err != nil {
		return nil, err
	}
	loadedDB := make([]apiservice.ServiceDef, 0)
	for _, file := range dbFiles {
		logger.Debugf("loading file '%s'", file)
		entries, err := apiservice.ServiceDefsFromFile(file)
		if err != nil {
			return nil, fmt.Errorf("loading file '%s': %v", file, err)
		}
		loadedDB = append(loadedDB, entries...)
	}
	return loadedDB, nil
}
